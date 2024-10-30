package chainio

//func TestAvsReader_AvsDirectory(t *testing.T) {
//	rpcUrl := "https://rpc.holesky.ethpandaops.io"
//	ethClient, err := eth.NewClient(rpcUrl)
//	if err != nil {
//		t.Errorf("Failed to create eth client: %v", err)
//		return
//	}
//
//	cfg := config.Config{
//		ServiceManagerAddr: common.HexToAddress("0x55301a1a9182732461567D0F9F34E3D9ce5F343A"),
//		EthHttpClient:      ethClient,
//	}
//
//	avsReader, err := NewAvsReader(cfg)
//	if err != nil {
//		t.Errorf("Failed to create AvsReader: %v", err)
//		return
//	}
//
//	addr, err := avsReader.AvsDirectory()
//	if err != nil {
//		t.Errorf("Failed to get AvsDirectory: %v", err)
//		return
//	}
//
//	t.Logf("AvsDirectory: %v", addr)
//
//}
