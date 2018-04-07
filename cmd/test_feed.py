import bitcoin
from ethereum import utils
from pycoin.serialize import b2h, h2b
import json


def generate_feed(block_number, nonce, ask_for_1000, bid_for_1000, private_key):
    hexX = h2b("%0.64X" % block_number)
    hexY = h2b("%0.64X" % nonce)
    hexZ = h2b("%0.64X" % ask_for_1000)
    hexW = h2b("%0.64X" % bid_for_1000)
    msg = hexX + hexY + hexZ + hexW
    msghash = b2h(utils.sha3(msg))
    V, R, S = bitcoin.ecdsa_raw_sign(msghash, private_key)
    # print("V R S")
    R = utils.int_to_hex(R)
    S = utils.int_to_hex(S)
    # print("0x%x" % V, R, S)

    dic = {}
    dic['status'] = "success"
    data = {}
    data['block_number'] = block_number
    data['nonce'] = nonce
    data['bid_for_1000'] = bid_for_1000
    data['ask_for_1000'] = ask_for_1000
    data['message'] = "0x" + b2h(msg)
    data['hash'] = "0x" + msghash
    data['signer'] = "0x" + b2h(utils.privtoaddr(private_key))
    data['v'] = V
    data['r'] = R
    data['s'] = S

    dic['data'] = data

    print(json.dumps(dic, indent=4, sort_keys=True))


private_key = bitcoin.sha256('some big long brainwallet password')


block_number = 5392391
nonce = 1523036543
ask_for_1000 = 48082
bid_for_1000 = 46440

generate_feed(block_number, nonce, ask_for_1000, bid_for_1000, private_key)
