# Validation plan


| Area | Test | Tool or command | Expected result | Evidence |
|---|---|---|---|---|
| Connectivity | Ping inside the same VLAN | ping | Successful ICMP echo replies | Screenshots / terminal logs |
| Connectivity | Ping between different VLANs and between sites | ping, traceroute | Successful routing between subnets and sites | Screenshots / terminal logs |
| Routing | Verify OSPF neighbors and routes | show ip ospf neighbor (FRRouting) | Routers R1 and R2 are full neighbors over wlan0 | Screenshots / terminal logs |
| Redundancy | Check EtherChannel status after disconnecting 1 of 4 cables | show etherchannel summary | Bundle stays up with remaining links | Screenshots |
| Redundancy | Verify RSTP status | show spanning-tree | Correct root bridge and port roles (forwarding/blocking) | Screenshots |
| Service access | Client from Site B logs in and chats with Server in Site A | Go client app | Successful login and message delivery | App screenshots / Wireshark captures |
| Security | Guest (VLAN 30) tries to reach servers (VLAN 20) | ping | Request times out (blocked by firewall) | Screenshots |
| Security | User (VLAN 10) tries to reach administration (VLAN 99) | ping, ssh | Request is denied (blocked by firewall) | Screenshots |
| Security | API request sent without a valid token | curl / Postman | HTTP 401 Unauthorized response | Screenshots / terminal logs |
| Security | Check database password storage | Database query | Only hashed passwords are stored, no plaintext | Screenshots |
| Security | Verify firewall ruleset | nft list ruleset | Expected firewall zones and rules are applied | Screenshots / logs |