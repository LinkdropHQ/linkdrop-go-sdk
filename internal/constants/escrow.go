package constants

import (
	"github.com/ethereum/go-ethereum/common"
)

var Escrows = map[string][]common.Address{
	"1": {
		common.HexToAddress("0x0522dd6e9f2beca1cd15a5fd275dc279a1a08eac"),
	},
	"2": {
		common.HexToAddress("0xad27383460183fd7e21b71df3b4cac9480eb9a75"),
		common.HexToAddress("0x0B79cC1E78C47fF08cA6f355e8aCD32AEa5bFe58"),
		common.HexToAddress("0xc4eb6e5933bc5e32dfd5c80baf143212a95549b3"),
	},
	"3": {
		common.HexToAddress("0x0b962bbbf101941d0d0ec1041d01668dac36647a"),
		common.HexToAddress("0x2d5dfe0e4582c905233df527242616017f36e192"),
		common.HexToAddress("0x021ccef76804c43da62b01652d41bcf6f6394731"),
	},
	"3.1": {
		common.HexToAddress("0x88d51990a3b962f975846f3688e36d2a1fc611f1"),
		common.HexToAddress("0x648b9a6c54890a8fb17de128c6352f621154f358"),
		common.HexToAddress("0x7143f68e689e8540a8eec26b482e1d4ac2e28794"),
		common.HexToAddress("0xe07fa88a10a915b7339aff050db82c0030bf6861"),
		common.HexToAddress("0x4366caf3963d147da4a4287061354058d871d1be"),
		common.HexToAddress("0x317d2501396fe75d997799bf3bdbc7cc6768b533"),
		common.HexToAddress("0x59548f7e4ef381df57a3e5dacbf2ab65111404d6"),
		common.HexToAddress("0xedfea6336c922f896c7e09ba282beb0cb4476675"),
		common.HexToAddress("0xff3471dfdc6f82694e5ad4d4e7ffedf23e1e38e0"),
		common.HexToAddress("0x139b79602b68e8198ea3d57f5e6311fd98262269"),
		common.HexToAddress("0xe0cec4f0b66257fc6b13652c303237de0fd92ed8"),
	},
	"3.2": {
		common.HexToAddress("0x5badb0143f69015c5c86cbd9373474a9c8ab713b"),
		common.HexToAddress("0x3c74782de03c0402d207fe41307fe50fe9b6b5c7"),
		common.HexToAddress("0xbe7b40eb3a9d85d3a76142cb637ab824f0d35ead"),
		common.HexToAddress("0x5fc1316119a1b7cec52a2984c62764343dca70c9"),
	},
}
