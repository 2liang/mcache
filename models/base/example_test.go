package base

func ExampleShopModel_1() {
	//DbShop主库
	DbShop.GetMaster()
	//DbShop从库
	DbShop.GetSlave()

	//Output:
	//
}
