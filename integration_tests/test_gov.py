from datetime import timedelta
from time import sleep

import pytest
from dateutil.parser import isoparse

from .utils import DEFAULT_BASE_PORT, astra_to_aastra, parse_events, wait_for_block, wait_for_block_time, wait_for_port, wait_for_new_blocks, GAS_USE

pytestmark = pytest.mark.gov


@pytest.mark.parametrize("vote_option", ["yes", "no"])
def test_param_proposal(astra, vote_option):
    fn_vote_option(astra, vote_option)


@pytest.mark.parametrize("vote_option", ["no_with_veto", "abstain", None])
def test_vote_option(astra, vote_option):
    fn_vote_option(astra, vote_option)


def fn_vote_option(astra, vote):
    """
    - send proposal to change max_validators
    - all validator vote same option (None means don't vote)
    - check the result
    - check deposit refunded
    """
    wait_for_new_blocks(astra.cosmos_cli(0), 1)
    wait_for_port(DEFAULT_BASE_PORT)
    max_validators = astra.cosmos_cli(0).staking_params()["max_validators"]

    rsp = astra.cosmos_cli(0).gov_propose(
        "team",
        "param-change",
        {
            "title": "Increase number of max validators",
            "description": "ditto",
            "changes": [
                {
                    "subspace": "staking",
                    "key": "MaxValidators",
                    "value": max_validators + 1,
                }
            ],
        },
    )
    assert rsp["code"] == 0, rsp["raw_log"]

    # get proposal_id
    ev = parse_events(rsp["logs"])["submit_proposal"]
    assert ev["proposal_type"] == "ParameterChange", rsp
    proposal_id = ev["proposal_id"]

    proposal = astra.cosmos_cli(0).query_proposal(proposal_id)
    assert proposal["content"]["changes"] == [
        {
            "subspace": "staking",
            "key": "MaxValidators",
            "value": str(max_validators + 1),
        }
    ], proposal
    assert proposal["status"] == "PROPOSAL_STATUS_DEPOSIT_PERIOD", proposal
    # deposit_amount >= gov:min_deposit
    deposit_amount = 10000000
    amount = astra.cosmos_cli(0).balance(astra.cosmos_cli(0).address("team"))
    print("team amount", amount)
    rsp = astra.cosmos_cli(0).gov_deposit("team", proposal_id, "%daastra" % deposit_amount)
    assert rsp["code"] == 0, rsp["raw_log"]
    assert astra.cosmos_cli(0).balance(astra.cosmos_cli(0).address("team")) == amount - deposit_amount

    proposal = astra.cosmos_cli(0).query_proposal(proposal_id)
    assert proposal["status"] == "PROPOSAL_STATUS_VOTING_PERIOD", proposal

    if vote is not None:
        wait_for_port(DEFAULT_BASE_PORT)
        # node #1
        rsp = astra.cosmos_cli(0).gov_vote("validator", proposal_id, vote)
        assert rsp["code"] == 0, rsp["raw_log"]
        wait_for_port(DEFAULT_BASE_PORT)
        # node #2
        rsp = astra.cosmos_cli(1).gov_vote("validator", proposal_id, vote)
        assert rsp["code"] == 0, rsp["raw_log"]
        wait_for_port(DEFAULT_BASE_PORT)
        assert (
                int(astra.cosmos_cli(1).query_tally(proposal_id)[vote])
                == astra.cosmos_cli(0).staking_pool()
        ), "all voted"
    else:
        wait_for_port(DEFAULT_BASE_PORT)
        assert astra.cosmos_cli(0).query_tally(proposal_id) == {
            "yes": "0",
            "no": "0",
            "abstain": "0",
            "no_with_veto": "0",
        }

    wait_for_block_time(
        astra.cosmos_cli(0), isoparse(proposal["voting_end_time"]) + timedelta(seconds=5)
    )
    wait_for_port(DEFAULT_BASE_PORT)

    proposal = astra.cosmos_cli(0).query_proposal(proposal_id)
    if vote == "yes":
        assert proposal["status"] == "PROPOSAL_STATUS_PASSED", proposal
    else:
        assert proposal["status"] == "PROPOSAL_STATUS_REJECTED", proposal

    new_max_validators = astra.cosmos_cli(0).staking_params()["max_validators"]
    if vote == "yes":
        assert new_max_validators == max_validators + 1
    else:
        assert new_max_validators == max_validators

    if vote in ("no_with_veto", None):
        # not refunded
        assert astra.cosmos_cli(0).balance(astra.cosmos_cli(0).address("team")) == amount - deposit_amount
    else:
        # refunded, no matter passed or rejected
        assert astra.cosmos_cli(0).balance(astra.cosmos_cli(0).address("team")) == amount


def test_deposit_period_expires(astra):
    """
    - proposal and partially deposit
    - wait for deposit period end and check
    - proposal deleted
    - no refund
    """
    wait_for_port(DEFAULT_BASE_PORT)
    amount1 = astra.cosmos_cli(0).balance(astra.cosmos_cli(0).address("team"))
    # deposit_amount < gov:min_deposit
    deposit_amount = 10000
    rsp = astra.cosmos_cli(0).gov_propose(
        "team",
        "param-change",
        {
            "title": "Increase number of max validators",
            "description": "ditto",
            "changes": [
                {
                    "subspace": "staking",
                    "key": "MaxValidators",
                    "value": 1,
                }
            ],
            "deposit": "%daastra" % deposit_amount,
        },
    )
    print(rsp)
    assert rsp["code"] == 0, rsp["raw_log"]
    ev = parse_events(rsp["logs"])["submit_proposal"]
    assert ev["proposal_type"] == "ParameterChange", rsp
    proposal_id = ev["proposal_id"]

    proposal = astra.cosmos_cli(0).query_proposal(proposal_id)
    assert proposal["total_deposit"] == [{"denom": "aastra", "amount": str(deposit_amount)}]

    assert astra.cosmos_cli(0).balance(astra.cosmos_cli(0).address("team")) == amount1 - deposit_amount - GAS_USE

    amount2 = astra.cosmos_cli(0).balance(astra.cosmos_cli(0).address("community"))
    rsp = astra.cosmos_cli(0).gov_deposit("community", proposal["proposal_id"], "%daastra" % deposit_amount)
    assert rsp["code"] == 0, rsp["raw_log"]
    proposal = astra.cosmos_cli(0).query_proposal(proposal_id)
    assert proposal["total_deposit"] == [{"denom": "aastra", "amount": str(deposit_amount * 2)}]

    assert astra.cosmos_cli(0).balance(astra.cosmos_cli(0).address("community")) == amount2 - deposit_amount

    # wait for deposit period passed
    wait_for_block_time(
        astra.cosmos_cli(0), isoparse(proposal["submit_time"]) + timedelta(seconds=15)
    )

    # proposal deleted
    with pytest.raises(Exception):
        proposal = astra.cosmos_cli(0).query_proposal(proposal_id)

    # deposits don't get refunded
    assert astra.cosmos_cli(0).balance(astra.cosmos_cli(0).address("team")) == amount1 - deposit_amount - GAS_USE
    assert astra.cosmos_cli(0).balance(astra.cosmos_cli(0).address("community")) == amount2 - deposit_amount


def test_community_pool_spend_proposal(astra):
    """
    - proposal a community pool spend
    - pass it
    """
    # need at least several blocks to populate community pool
    wait_for_port(DEFAULT_BASE_PORT)
    wait_for_new_blocks(astra.cosmos_cli(0), 1)
    amount = int(astra.cosmos_cli(0).distribution_community())
    print("Distribution community amount: %d" % amount)
    assert amount > 0, "need positive pool to proceed this test"

    recipient = astra.cosmos_cli(0).address("community")
    old_amount = astra.cosmos_cli(0).balance(recipient)
    print("Community old amount: %d" % old_amount)

    deposit_amount = 10000001

    rsp = astra.cosmos_cli(0).gov_propose(
        "community",
        "community-pool-spend",
        {
            "title": "Community Pool Spend",
            "description": "Pay me some astra!",
            "recipient": recipient,
            "amount": "%daastra" % amount,
            "deposit": "%daastra" % deposit_amount,
        },
    )
    assert rsp["code"] == 0, rsp["raw_log"]

    # get proposal_id
    ev = parse_events(rsp["logs"])["submit_proposal"]
    assert ev["proposal_type"] == "CommunityPoolSpend", rsp
    proposal_id = ev["proposal_id"]

    # vote
    rsp = astra.cosmos_cli(0).gov_vote("validator", proposal_id, "yes")
    assert rsp["code"] == 0, rsp["raw_log"]
    rsp = astra.cosmos_cli(1).gov_vote("validator", proposal_id, "yes")
    assert rsp["code"] == 0, rsp["raw_log"]

    # wait for voting period end
    proposal = astra.cosmos_cli(0).query_proposal(proposal_id)
    assert proposal["total_deposit"] == [{"denom": "aastra", "amount": str(deposit_amount)}]
    assert proposal["status"] == "PROPOSAL_STATUS_VOTING_PERIOD", proposal
    wait_for_block_time(
        astra.cosmos_cli(0), isoparse(proposal["voting_end_time"]) + timedelta(seconds=1)
    )

    proposal = astra.cosmos_cli(0).query_proposal(proposal_id)
    assert proposal["status"] == "PROPOSAL_STATUS_PASSED", proposal
    assert astra.cosmos_cli(0).balance(recipient) == old_amount + amount - GAS_USE


def test_change_vote(astra):
    """
    - submit proposal with deposit
    - vote yes
    - check tally
    - change vote
    - check tally
    """
    wait_for_port(DEFAULT_BASE_PORT)
    wait_for_new_blocks(astra.cosmos_cli(0), 1)
    deposit_amount = 10000000
    rsp = astra.cosmos_cli(0).gov_propose(
        "community",
        "param-change",
        {
            "title": "Increase number of max validators",
            "description": "ditto",
            "changes": [
                {
                    "subspace": "staking",
                    "key": "MaxValidators",
                    "value": 1,
                }
            ],
            "deposit": "%dastra" % deposit_amount,
        },
    )
    assert rsp["code"] == 0, rsp["raw_log"]

    voting_power = int(
        astra.cosmos_cli(0).validator(astra.cosmos_cli(0).address("validator", bech="val"))["tokens"]
    )
    print(voting_power)

    proposal_id = parse_events(rsp["logs"])["submit_proposal"]["proposal_id"]

    rsp = astra.cosmos_cli(0).gov_vote("validator", proposal_id, "yes")
    assert rsp["code"] == 0, rsp["raw_log"]

    # change vote to no
    rsp = astra.cosmos_cli(0).gov_vote("validator", proposal_id, "no")
    assert rsp["code"] == 0, rsp["raw_log"]


def test_inherit_vote(astra):
    """
    - submit proposal with deposits
    - A delegate to V
    - V vote Yes
    - check tally: {yes: a + v}
    - A vote No
    - change tally: {yes: v, no: a}
    """
    wait_for_port(DEFAULT_BASE_PORT)
    wait_for_new_blocks(astra.cosmos_cli(0), 1)
    deposit_amount = 10000000
    rsp = astra.cosmos_cli(0).gov_propose(
        "community",
        "param-change",
        {
            "title": "Increase number of max validators",
            "description": "ditto",
            "changes": [
                {
                    "subspace": "staking",
                    "key": "MaxValidators",
                    "value": 1,
                }
            ],
            "deposit": "%daastra" % deposit_amount,
        },
    )
    assert rsp["code"] == 0, rsp["raw_log"]
    proposal_id = parse_events(rsp["logs"])["submit_proposal"]["proposal_id"]

    delegate_amount = 10

    staked_amount_val1 = 1001

    voter1 = astra.cosmos_cli(0).address("community")
    astra.cosmos_cli(0).delegate_amount(
        # to_addr       amount      from_addr
        astra.cosmos_cli(0).address("validator", bech="val"), "%daastra" % delegate_amount, voter1  # delegate to validator #2
    )

    rsp = astra.cosmos_cli(0).gov_vote("validator", proposal_id, "yes")
    assert rsp["code"] == 0, rsp["raw_log"]
    wait_for_new_blocks(astra.cosmos_cli(0), 1)
    assert astra.cosmos_cli(0).query_tally(proposal_id) == {
        "yes": str(astra_to_aastra(staked_amount_val1) + delegate_amount),
        "no": "0",
        "abstain": "0",
        "no_with_veto": "0",
    }

    rsp = astra.cosmos_cli(0).gov_vote(voter1, proposal_id, "no")
    assert rsp["code"] == 0, rsp["raw_log"]

    assert astra.cosmos_cli(0).query_tally(proposal_id) == {
        "yes": str(astra_to_aastra(staked_amount_val1)),
        "no": str(delegate_amount),
        "abstain": "0",
        "no_with_veto": "0",
    }
