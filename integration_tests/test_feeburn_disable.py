from time import sleep
import pytest
from pathlib import Path

from eth_bloom import BloomFilter
from eth_utils import abi, big_endian_to_int
from hexbytes import HexBytes
from integration_tests.network import setup_astra

from integration_tests.utils import astra_to_aastra, deploy_contract, CONTRACTS, KEYS, ADDRS, send_transaction, wait_for_block, wait_for_new_blocks, GAS_USE, DEFAULT_BASE_PORT

pytestmark = pytest.mark.feeburn


@pytest.fixture(scope="module")
def astra(tmp_path_factory):
    path = tmp_path_factory.mktemp("astra")
    cfg = Path(__file__).parent / "configs/feeburn_disable.yaml"
    yield from setup_astra(path, DEFAULT_BASE_PORT, cfg)


def test_transfer(astra):
    """
    check simple transfer tx success
    - send 1astra from team to treasury
    """
    team_addr = astra.cosmos_cli(0).address("team")
    addr = "astra1wyzq5uv53tf7lqxn4qjmujlg8fcsmtafhs97ph"

    amount_astra = 1
    amount_aastra = astra_to_aastra(amount_astra)
    fee_coins = 1000000000
    old_block_height =  astra.cosmos_cli(0).block_height()
    print("old_block_height", old_block_height)
    old_block_provisions = round(astra.cosmos_cli(0).annual_provisions() / 6311520)
    print("old_block_provisions", old_block_provisions)
    old_total_supply = int(float(astra.cosmos_cli(0).total_supply()["supply"][0]["amount"]))
    print("old_total_supply", old_total_supply)
    tx = astra.cosmos_cli(0).transfer(team_addr, addr, str(amount_astra) + "astra", fees="%saastra" % fee_coins)
    tx_block_height = int(tx["height"])
    print("tx_block_height", tx_block_height)
    assert tx["logs"] == [
        {
            "events": [
                {
                    "attributes": [
                        {"key": "receiver", "value": addr},
                        {"key": "amount", "value": str(amount_aastra) + "aastra"},
                    ],
                    "type": "coin_received",
                },
                {
                    "attributes": [
                        {"key": "spender", "value": team_addr},
                        {"key": "amount", "value": str(amount_aastra) + "aastra"},
                    ],
                    "type": "coin_spent",
                },
                {
                    "attributes": [
                        {"key": "action", "value": "/cosmos.bank.v1beta1.MsgSend"},
                        {"key": "sender", "value": team_addr},
                        {"key": "module", "value": "bank"},
                    ],
                    "type": "message",
                },
                {
                    "attributes": [
                        {"key": "recipient", "value": addr},
                        {"key": "sender", "value": team_addr},
                        {"key": "amount", "value": str(amount_aastra) + "aastra"},
                    ],
                    "type": "transfer",
                },
            ],
            "log": "",
            "msg_index": 0,
        }
    ]
    new_block_provisions = round(astra.cosmos_cli(0).annual_provisions() / 6311520)
    print("new_block_provisions", new_block_provisions)
    # wait_for_new_blocks(astra.cosmos_cli(0), 1)
    # block_provisions1 = int(int(astra.cosmos_cli(0).annual_provisions()) / 6311520)
    # print("block_provisions1", block_provisions1)
    new_total_supply = int(float(astra.cosmos_cli(0).total_supply()["supply"][0]["amount"]))
    if tx_block_height - old_block_height == 2:
        diff_balance = new_block_provisions + old_block_provisions - (new_total_supply - old_total_supply)
    else:
        diff_balance = new_block_provisions - (new_total_supply - old_total_supply)
    print(new_block_provisions, new_total_supply - old_total_supply, diff_balance)
    print("block_height", astra.cosmos_cli(0).block_height())
    assert diff_balance < 0