package schedulor

//func TestMyListener(t *testing.T) {
//	if testing.Short() {
//		t.Skip("skipping test in short mode.")
//	}
//	logrus.SetLevel(logrus.DebugLevel)
//	config.InitConfig("../config/basic_strategy_config.json", "../config/basic_config.json", "../config/secret_config.json")
//	client, err := ethclient.Dial("wss://optimism-sepolia.infura.io/ws/v3/0284f5a9fc55476698079b24e2f97909")
//	if err != nil {
//		panic(err)
//	}
//	listener := EventListener{
//		logCh:               make(chan types.Log, 10),
//		entryPointAddresses: []common.Address{common.HexToAddress("0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789")},
//		paymasterAddresses:  []common.Address{common.HexToAddress("0x8817340e0a3435E06254f2ed411E6418cd070D6F")},
//		client:              client,
//		network:             global_const.OptimismSepolia,
//	}
//
//	listener.Listen()
//}
//
//func TestAddListener(t *testing.T) {
//	if testing.Short() {
//		t.Skip("skipping test in short mode.")
//	}
//	client, err := ethclient.Dial("wss://optimism-sepolia.infura.io/ws/v3/0284f5a9fc55476698079b24e2f97909")
//	if err != nil {
//		panic(err)
//	}
//	chainId, err := client.ChainID(context.Background())
//	if err != nil {
//		panic(err)
//	}
//	t.Logf("chainId: %v", chainId.String())
//	entryPointAdd := common.HexToAddress("0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789")
//	//topics := make([][]interface{}, 0)
//	eventId := userOpEvent.ID
//	var eventIds []interface{}
//	eventIds = append(eventIds, eventId)
//	querys := make([][]interface{}, 0)
//	querys = append(querys, eventIds)
//
//	var userOpHashRule []interface{}
//	querys = append(querys, userOpHashRule)
//	var senderRule []interface{}
//	querys = append(querys, senderRule)
//	var paymasterRule []interface{}
//	paymasterRule = append(paymasterRule, common.HexToAddress("0x8817340e0a3435E06254f2ed411E6418cd070D6F"))
//	querys = append(querys, paymasterRule)
//
//	//t.Log("event ", eventId.String())
//	//topics = append(topics, []common.Hash{eventId})
//	//payAddress := common.HexToAddress("0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789")
//	//topics[3] = []common.Hash{payAddress.Hash()}
//
//	topics, err := abi.MakeTopics(querys...)
//	query := ethereum.FilterQuery{
//		Addresses: []common.Address{entryPointAdd},
//		Topics:    topics,
//	}
//
//	logs := make(chan types.Log, 11)
//
//	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
//	if err != nil {
//		panic(err)
//	}
//
//	for {
//		select {
//		case err := <-sub.Err():
//			t.Fatal(err)
//		case vLog := <-logs:
//			fmt.Println("tx Hash: ", vLog.TxHash.String())
//			if len(vLog.Topics) == 0 {
//				t.Logf("vLog.Topics is empty")
//				t.Fatal(err)
//			}
//			if vLog.Topics[0] != eventId {
//				t.Logf("vLog.Topics[0] != eventId")
//				t.Fatal(err)
//			}
//			if len(vLog.Topics) < 4 {
//				t.Logf("vLog.Topics length < 4")
//				t.Fatal(err)
//			}
//			if len(vLog.Data) > 0 {
//				eventObj := &ContractUserOperationEvent{}
//				err := entryPointABI.UnpackIntoInterface(eventObj, "UserOperationEvent", vLog.Data)
//				t.Logf("Data: %v", hex.EncodeToString(vLog.Data))
//				if err != nil {
//					t.Errorf("UnpackIntoInterface failed: %v", err)
//					t.Fatal(err)
//				}
//				eventObj.UserOpHash = common.BytesToHash(vLog.Topics[1].Bytes())
//				eventObj.Sender = common.BytesToAddress(vLog.Topics[2].Bytes())
//				eventObj.Paymaster = common.BytesToAddress(vLog.Topics[3].Bytes())
//
//				jsonStr, _ := json.Marshal(eventObj)
//				t.Logf("userOpEvent: %v", string(jsonStr))
//			}
//		}
//	}
//}
//
////	func TestAddListener2(t *testing.T) {
////		config.InitConfig("../config/basic_strategy_config.json", "../config/basic_config.json", "../config/secret_config.json")
////
////		add := common.HexToAddress("0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789")
////		entryPpointContract, err := executor.GetEntryPoint06(&add)
////		if err != nil {
////			t.Fatal(err)
////		}
////		it, err := entryPpointContract.FilterUserOperationEvent(&bind.FilterOpts{}, nil, nil, nil)
////		if err != nil {
////			t.Fatal(err)
////		}
////		for {
////			if it.Next() {
////				fmt.Println(it.Event)
////			}
////		}
////	}
//func TestABi(t *testing.T) {
//	// 0x
//	abi, _ := contract_entrypoint_v06.ContractMetaData.GetAbi()
//	event := abi.Events["UserOperationEvent"]
//	id := event.ID
//	t.Log(id)
//
//	abi07, _ := contract_entrypoint_v07.ContractMetaData.GetAbi()
//	event07 := abi07.Events["UserOperationEvent"]
//	id07 := event07.ID
//	t.Log(id07)
//}
