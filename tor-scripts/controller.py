# %%
# 包引入
import requests
from stem.descriptor import DocumentHandler, parse_file
from stem.descriptor import parse_file, DocumentHandler
from stem.control import EventType
from stem.control import Controller

# %%
# 连接到 Tor 控制器

c = Controller.from_port(port=9051)
c.authenticate("my_password")

# %%
# 初始化 IP 获取


def getIP(proxy=None):
    proxies = {}
    if proxy != None:
        proxies["http"] = proxy
        proxies["https"] = proxy
    data = requests.get("http://ip-api.com/json", proxies=proxies).json()
    if data["status"] == "success":
        print(data["query"], data["city"])


getIP()
getIP("socks5://127.0.0.1:9050")

# %%
# 节点信息获取

fields = [
    "fingerprint", "or_addresses", "identifier_type",
    "identifier", "digest", "bandwidth",
    "measured", "is_unmeasured", "unrecognized_bandwidth_entries",
    "exit_policy", "protocols", "microdescriptor_hashes",
    "document", "nickname", "published",
    "address", "or_port", "dir_port",
    "flags", "version", "version_line",
]
relays = {}
relays_nickname = {}
for relay in c.get_network_statuses():
    relay_info = {
        f: getattr(relay, f)
        for f in fields
    }
    relays_nickname[relay.nickname] = relay_info
    relays[str(relay.fingerprint)] = relay_info

# %%
# 中继信息测试

print(relays["0A5476625BAC0573A82FCB432ACED70E1ED83E93"])
print(relays_nickname["OhYee"])
print(len(relays))

# %%
# 自主构建电路

circuit_id = c.new_circuit(
    [
        # "B7D773E2196DB5B4129D3C050C063DB0F8DF91A5",
        "0A5476625BAC0573A82FCB432ACED70E1ED83E93",
        # "8CA736F332937E4FA0F5F48488B823F31883123C",
    ],
    await_build=True
)
print(circuit_id)


def attach_stream(stream):
    if stream.status == 'NEW':
        c.attach_stream(stream.id, circuit_id)


c.add_event_listener(attach_stream, EventType.STREAM)
try:
    c.set_conf('__LeaveStreamsUnattached', '1')
    getIP("socks5://127.0.0.1:9050")
finally:
    c.remove_event_listener(attach_stream)
    c.reset_conf('__LeaveStreamsUnattached')

# %%
# 获取共识信息

consensus = {}
consensus_nickname = {}

with open('/home/ohyee/.tor/cached-consensus', 'rb') as consensus_file:
    routers = parse_file(consensus_file, 'network-status-consensus-3 1.0',
                         document_handler=DocumentHandler.ENTRIES)
    for router in routers:
        consensus_nickname[router.nickname] = router
        consensus[router.fingerprint] = router


# %%
print(len(consensus))


# %%
# 对比共识和节点内容

cnt = 0
for relay in relays.values():
    temp = consensus_nickname.get(relay["nickname"], None)
    if temp == None:
        cnt += 1
        print(relay["nickname"], relay["fingerprint"])
print(cnt)

# %%
# 查看当前建立的电路列表

circuits = []
for circuit in c.get_circuits():
    circuits.append(circuit.id)
    print(circuit.id)
    for relay in circuit.path:
        print(" ", relay)


# %%
# 关闭所有电路

for cid in circuits:
    c.close_circuit(cid)


# %%
# 关闭控制连接

c.close()
