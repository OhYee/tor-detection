# %%
# 初始化

import requests
from stem.descriptor import DocumentHandler, parse_file
from stem.descriptor import parse_file, DocumentHandler
from stem.control import EventType
from stem.control import Controller

c = None


def getIP(proxy=None):
    proxies = {}
    if proxy != None:
        proxies["http"] = proxy
        proxies["https"] = proxy
    data = requests.get("http://ip-api.com/json", proxies=proxies).json()
    if data["status"] == "success":
        print(data["query"], data["city"])


relays = {}
relays_nickname = {}


def initial_relays():
    global c, relays, relays_nickname
    fields = [
        "fingerprint", "or_addresses", "identifier_type",
        "identifier", "digest", "bandwidth",
        "measured", "is_unmeasured", "unrecognized_bandwidth_entries",
        "exit_policy", "protocols", "microdescriptor_hashes",
        "document", "nickname", "published",
        "address", "or_port", "dir_port",
        "flags", "version", "version_line",
    ]
    for relay in c.get_network_statuses():
        relay_info = {
            f: getattr(relay, f)
            for f in fields
        }
        relays_nickname[relay.nickname] = relay_info
        relays[str(relay.fingerprint)] = relay_info


# 自主链路 id
my_circuit_id = -1


def attach_stream(stream):
    '''
    自主链路绑定
    '''
    global c, my_circuit_id
    print(stream)
    if stream.status == 'NEW' and my_circuit_id != -1:
        c.attach_stream(stream.id, my_circuit_id)


def use_my_circuit(circuit_id):
    global c, my_circuit_id
    if my_circuit_id == -1:
        my_circuit_id = circuit_id
        c.add_event_listener(attach_stream, EventType.STREAM)
        c.set_conf('__LeaveStreamsUnattached', '1')
    else:
        print("已经设置过自主电路，可能需要重新配置")


def close_my_circuit():
    global c, my_circuit_id
    if my_circuit_id != -1:
        c.remove_event_listener(attach_stream)
        c.reset_conf('__LeaveStreamsUnattached')
        c.close_circuit(my_circuit_id)
        my_circuit_id = -1


consensus = {}
consensus_nickname = {}


def initial_consensus():
    global c, consensus, consensus_nickname
    with open('/home/ohyee/.tor/cached-consensus', 'rb') as consensus_file:
        routers = parse_file(consensus_file, 'network-status-consensus-3 1.0',
                             document_handler=DocumentHandler.ENTRIES)
        for router in routers:
            consensus_nickname[router.nickname] = router
            consensus[router.fingerprint] = router


def get_all_circuits():
    '''
    获取所有电路
    '''
    global c
    circuits = []
    for circuit in c.get_circuits():
        circuits.append(circuit)
        print(circuit.id, circuit.purpose, circuit.status, circuit.build_flags)
        for relay in circuit.path:
            print(" ", relay)
    return circuits


def show_circuits():
    '''
    查看所有电路
    '''
    circuits = get_all_circuits()
    for circuit in circuits:
        print(circuit.id, circuit.purpose, circuit.status, circuit.build_flags)
        for relay in circuit.path:
            print(" ", relay)


def close_circuits():
    '''
    删除所有电路
    '''
    circuits = get_all_circuits()
    for circuit in circuits:
        c.close_circuit(circuit.id)


def new_circuit(path=[
    # "B7D773E2196DB5B4129D3C050C063DB0F8DF91A5",
    "0A5476625BAC0573A82FCB432ACED70E1ED83E93",
    "8CA736F332937E4FA0F5F48488B823F31883123C",
]):
    '''
    建立自主电路
    '''
    global c
    circuit_id = c.new_circuit(
        path=path,
        await_build=True
    )
    return circuit_id


def initial():
    '''
    初始化数据
    '''
    global c
    if c != None:
        close()
    c = Controller.from_port(port=9051)
    c.authenticate("my_password")

    initial_consensus()
    initial_relays()

    getIP()
    getIP("socks5://127.0.0.1:9050")


def close():
    global c

    c.close()


initial()

# %%
# 结束
close()

#######################################

# %%
# 获取所有电路
show_circuits()

# %%
# 关闭所有电路
close_circuits()

# %%
# 使用自主链路
use_my_circuit(new_circuit(
    path=[
        # "B7D773E2196DB5B4129D3C050C063DB0F8DF91A5",
        "0A5476625BAC0573A82FCB432ACED70E1ED83E93",
        "8CA736F332937E4FA0F5F48488B823F31883123C",
    ]
))


# %%
# 获取 IP
getIP("socks5://127.0.0.1:9050")

# %%
# 不使用自主链路
close_my_circuit()

# %%
