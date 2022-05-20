package evm

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/33cn/chain33-sdk-go/client"
	"github.com/33cn/chain33-sdk-go/crypto"
	"github.com/33cn/chain33-sdk-go/types"
	"github.com/stretchr/testify/assert"
)

var (
	deployAddress    = "14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
	deployPrivateKey = "CC38546E9E659D15E6B4893F0AB32A06D103931A8230B0BDE71459D2B27D6944"

	withholdAddress    = "14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
	withholdPrivateKey = "CC38546E9E659D15E6B4893F0AB32A06D103931A8230B0BDE71459D2B27D6944"

	useraAddress    = "17RH6oiMbUjat3AAyQeifNiACPFefvz3Au"
	useraPrivateKey = "56d1272fcf806c3c5105f3536e39c8b33f88cb8971011dfe5886159201884763"

	url      = "http://127.0.0.1:8901"
	paraName = "user.p.mbaas."
	//paraName = ""

	// solidity合约源码见：./solidity/ERC1155.sol
	codes = "60806040523480156200001157600080fd5b506040518060200160405280600081525062000033816200007b60201b60201c565b5033600360006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550620001ac565b80600290805190602001906200009392919062000097565b5050565b828054620000a59062000147565b90600052602060002090601f016020900481019282620000c9576000855562000115565b82601f10620000e457805160ff191683800117855562000115565b8280016001018555821562000115579182015b8281111562000114578251825591602001919060010190620000f7565b5b50905062000124919062000128565b5090565b5b808211156200014357600081600090555060010162000129565b5090565b600060028204905060018216806200016057607f821691505b602082108114156200017757620001766200017d565b5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b612fe880620001bc6000396000f3fe608060405234801561001057600080fd5b50600436106100a85760003560e01c80634e1273f4116100715780634e1273f414610175578063a22cb465146101a5578063ab918735146101c1578063b2bdfa7b146101dd578063e985e9c5146101fb578063f242432a1461022b576100a8565b8062fdd58e146100ad57806301ffc9a7146100dd5780630e89341c1461010d5780631f04bcc71461013d5780632eb2c2d614610159575b600080fd5b6100c760048036038101906100c29190611f2e565b610247565b6040516100d4919061269e565b60405180910390f35b6100f760048036038101906100f29190612025565b610310565b60405161010491906124c1565b60405180910390f35b61012760048036038101906101229190612077565b6103f2565b60405161013491906124dc565b60405180910390f35b61015760048036038101906101529190611e47565b610497565b005b610173600480360381019061016e9190611cf9565b6105f5565b005b61018f600480360381019061018a9190611fb9565b610696565b60405161019c9190612468565b60405180910390f35b6101bf60048036038101906101ba9190611ef2565b610847565b005b6101db60048036038101906101d69190611f6a565b61085d565b005b6101e561087e565b6040516101f2919061238b565b60405180910390f35b61021560048036038101906102109190611cbd565b6108a4565b60405161022291906124c1565b60405180910390f35b61024560048036038101906102409190611db8565b610938565b005b60008073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1614156102b8576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016102af9061257e565b60405180910390fd5b60008083815260200190815260200160002060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b60007fd9b67a26000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff191614806103db57507f0e89341c000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916145b806103eb57506103ea826109d9565b5b9050919050565b60606004600083815260200190815260200160002080546104129061296a565b80601f016020809104026020016040519081016040528092919081815260200182805461043e9061296a565b801561048b5780601f106104605761010080835404028352916020019161048b565b820191906000526020600020905b81548152906001019060200180831161046e57829003601f168201915b50505050509050919050565b81518351146104db576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016104d29061253e565b60405180910390fd5b805183511461051f576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610516906125fe565b60405180910390fd5b61053a84848460405180602001604052806000815250610a43565b6000815111156105ef5760005b83518110156105ed576105da84828151811061058c577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60200260200101518383815181106105cd577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b6020026020010151610cbc565b80806105e5906129cd565b915050610547565b505b50505050565b6105fd610ce8565b73ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff16148061064357506106428561063d610ce8565b6108a4565b5b610682576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610679906125be565b60405180910390fd5b61068f8585858585610cf0565b5050505050565b606081518351146106dc576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106d39061263e565b60405180910390fd5b6000835167ffffffffffffffff81111561071f577f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b60405190808252806020026020018201604052801561074d5781602001602082028036833780820191505090505b50905060005b845181101561083c576107e6858281518110610798577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60200260200101518583815181106107d9577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b6020026020010151610247565b82828151811061081f577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60200260200101818152505080610835906129cd565b9050610753565b508091505092915050565b610859610852610ce8565b838361105e565b5050565b6108793384848460405180602001604052806000815250610938565b505050565b600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16905092915050565b610940610ce8565b73ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff161480610986575061098585610980610ce8565b6108a4565b5b6109c5576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016109bc9061255e565b60405180910390fd5b6109d285858585856111cb565b5050505050565b60007f01ffc9a7000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916827bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916149050919050565b600073ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff161415610ab3576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610aaa9061267e565b60405180910390fd5b8151835114610af7576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610aee9061265e565b60405180910390fd5b6000610b01610ce8565b9050610b1281600087878787611467565b60005b8451811015610c1757838181518110610b57577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b6020026020010151600080878481518110610b9b577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b6020026020010151815260200190815260200160002060008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254610bfd919061285e565b925050819055508080610c0f906129cd565b915050610b15565b508473ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff167f4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb8787604051610c8f92919061248a565b60405180910390a4610ca68160008787878761146f565b610cb581600087878787611477565b5050505050565b80600460008481526020019081526020016000209080519060200190610ce392919061192e565b505050565b600033905090565b8151835114610d34576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610d2b9061265e565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff161415610da4576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610d9b9061259e565b60405180910390fd5b6000610dae610ce8565b9050610dbe818787878787611467565b60005b8451811015610fbb576000858281518110610e05577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b602002602001015190506000858381518110610e4a577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60200260200101519050600080600084815260200190815260200160002060008b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905081811015610eeb576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610ee2906125de565b60405180910390fd5b81810360008085815260200190815260200160002060008c73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508160008085815260200190815260200160002060008b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254610fa0919061285e565b9250508190555050505080610fb4906129cd565b9050610dc1565b508473ffffffffffffffffffffffffffffffffffffffff168673ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff167f4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb878760405161103292919061248a565b60405180910390a461104881878787878761146f565b611056818787878787611477565b505050505050565b8173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1614156110cd576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016110c49061261e565b60405180910390fd5b80600160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff0219169083151502179055508173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167f17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31836040516111be91906124c1565b60405180910390a3505050565b600073ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff16141561123b576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016112329061259e565b60405180910390fd5b6000611245610ce8565b905060006112528561165e565b9050600061125f8561165e565b905061126f838989858589611467565b600080600088815260200190815260200160002060008a73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905085811015611306576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016112fd906125de565b60405180910390fd5b85810360008089815260200190815260200160002060008b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508560008089815260200190815260200160002060008a73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546113bb919061285e565b925050819055508773ffffffffffffffffffffffffffffffffffffffff168973ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff167fc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f628a8a6040516114389291906126b9565b60405180910390a461144e848a8a86868a61146f565b61145c848a8a8a8a8a611724565b505050505050505050565b505050505050565b505050505050565b6114968473ffffffffffffffffffffffffffffffffffffffff1661190b565b15611656578373ffffffffffffffffffffffffffffffffffffffff1663bc197c8187878686866040518663ffffffff1660e01b81526004016114dc9594939291906123a6565b602060405180830381600087803b1580156114f657600080fd5b505af192505050801561152757506040513d601f19601f82011682018060405250810190611524919061204e565b60015b6115cd57611533612aa3565b806308c379a014156115905750611548612ec0565b806115535750611592565b806040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161158791906124dc565b60405180910390fd5b505b6040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016115c4906124fe565b60405180910390fd5b63bc197c8160e01b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916817bffffffffffffffffffffffffffffffffffffffffffffffffffffffff191614611654576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161164b9061251e565b60405180910390fd5b505b505050505050565b60606000600167ffffffffffffffff8111156116a3577f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040519080825280602002602001820160405280156116d15781602001602082028036833780820191505090505b509050828160008151811061170f577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60200260200101818152505080915050919050565b6117438473ffffffffffffffffffffffffffffffffffffffff1661190b565b15611903578373ffffffffffffffffffffffffffffffffffffffff1663f23a6e6187878686866040518663ffffffff1660e01b815260040161178995949392919061240e565b602060405180830381600087803b1580156117a357600080fd5b505af19250505080156117d457506040513d601f19601f820116820180604052508101906117d1919061204e565b60015b61187a576117e0612aa3565b806308c379a0141561183d57506117f5612ec0565b80611800575061183f565b806040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161183491906124dc565b60405180910390fd5b505b6040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611871906124fe565b60405180910390fd5b63f23a6e6160e01b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916817bffffffffffffffffffffffffffffffffffffffffffffffffffffffff191614611901576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016118f89061251e565b60405180910390fd5b505b505050505050565b6000808273ffffffffffffffffffffffffffffffffffffffff163b119050919050565b82805461193a9061296a565b90600052602060002090601f01602090048101928261195c57600085556119a3565b82601f1061197557805160ff19168380011785556119a3565b828001600101855582156119a3579182015b828111156119a2578251825591602001919060010190611987565b5b5090506119b091906119b4565b5090565b5b808211156119cd5760008160009055506001016119b5565b5090565b60006119e46119df84612707565b6126e2565b90508083825260208201905082856020860282011115611a0357600080fd5b60005b85811015611a335781611a198882611b82565b845260208401935060208301925050600181019050611a06565b5050509392505050565b6000611a50611a4b84612733565b6126e2565b9050808382526020820190508260005b85811015611a905781358501611a768882611c7e565b845260208401935060208301925050600181019050611a60565b5050509392505050565b6000611aad611aa88461275f565b6126e2565b90508083825260208201905082856020860282011115611acc57600080fd5b60005b85811015611afc5781611ae28882611ca8565b845260208401935060208301925050600181019050611acf565b5050509392505050565b6000611b19611b148461278b565b6126e2565b905082815260208101848484011115611b3157600080fd5b611b3c848285612928565b509392505050565b6000611b57611b52846127bc565b6126e2565b905082815260208101848484011115611b6f57600080fd5b611b7a848285612928565b509392505050565b600081359050611b9181612f56565b92915050565b600082601f830112611ba857600080fd5b8135611bb88482602086016119d1565b91505092915050565b600082601f830112611bd257600080fd5b8135611be2848260208601611a3d565b91505092915050565b600082601f830112611bfc57600080fd5b8135611c0c848260208601611a9a565b91505092915050565b600081359050611c2481612f6d565b92915050565b600081359050611c3981612f84565b92915050565b600081519050611c4e81612f84565b92915050565b600082601f830112611c6557600080fd5b8135611c75848260208601611b06565b91505092915050565b600082601f830112611c8f57600080fd5b8135611c9f848260208601611b44565b91505092915050565b600081359050611cb781612f9b565b92915050565b60008060408385031215611cd057600080fd5b6000611cde85828601611b82565b9250506020611cef85828601611b82565b9150509250929050565b600080600080600060a08688031215611d1157600080fd5b6000611d1f88828901611b82565b9550506020611d3088828901611b82565b945050604086013567ffffffffffffffff811115611d4d57600080fd5b611d5988828901611beb565b935050606086013567ffffffffffffffff811115611d7657600080fd5b611d8288828901611beb565b925050608086013567ffffffffffffffff811115611d9f57600080fd5b611dab88828901611c54565b9150509295509295909350565b600080600080600060a08688031215611dd057600080fd5b6000611dde88828901611b82565b9550506020611def88828901611b82565b9450506040611e0088828901611ca8565b9350506060611e1188828901611ca8565b925050608086013567ffffffffffffffff811115611e2e57600080fd5b611e3a88828901611c54565b9150509295509295909350565b60008060008060808587031215611e5d57600080fd5b6000611e6b87828801611b82565b945050602085013567ffffffffffffffff811115611e8857600080fd5b611e9487828801611beb565b935050604085013567ffffffffffffffff811115611eb157600080fd5b611ebd87828801611beb565b925050606085013567ffffffffffffffff811115611eda57600080fd5b611ee687828801611bc1565b91505092959194509250565b60008060408385031215611f0557600080fd5b6000611f1385828601611b82565b9250506020611f2485828601611c15565b9150509250929050565b60008060408385031215611f4157600080fd5b6000611f4f85828601611b82565b9250506020611f6085828601611ca8565b9150509250929050565b600080600060608486031215611f7f57600080fd5b6000611f8d86828701611b82565b9350506020611f9e86828701611ca8565b9250506040611faf86828701611ca8565b9150509250925092565b60008060408385031215611fcc57600080fd5b600083013567ffffffffffffffff811115611fe657600080fd5b611ff285828601611b97565b925050602083013567ffffffffffffffff81111561200f57600080fd5b61201b85828601611beb565b9150509250929050565b60006020828403121561203757600080fd5b600061204584828501611c2a565b91505092915050565b60006020828403121561206057600080fd5b600061206e84828501611c3f565b91505092915050565b60006020828403121561208957600080fd5b600061209784828501611ca8565b91505092915050565b60006120ac838361236d565b60208301905092915050565b6120c1816128b4565b82525050565b60006120d2826127fd565b6120dc818561282b565b93506120e7836127ed565b8060005b838110156121185781516120ff88826120a0565b975061210a8361281e565b9250506001810190506120eb565b5085935050505092915050565b61212e816128c6565b82525050565b600061213f82612808565b612149818561283c565b9350612159818560208601612937565b61216281612ac5565b840191505092915050565b600061217882612813565b612182818561284d565b9350612192818560208601612937565b61219b81612ac5565b840191505092915050565b60006121b360348361284d565b91506121be82612ae3565b604082019050919050565b60006121d660288361284d565b91506121e182612b32565b604082019050919050565b60006121f960218361284d565b915061220482612b81565b604082019050919050565b600061221c60298361284d565b915061222782612bd0565b604082019050919050565b600061223f602a8361284d565b915061224a82612c1f565b604082019050919050565b600061226260258361284d565b915061226d82612c6e565b604082019050919050565b600061228560328361284d565b915061229082612cbd565b604082019050919050565b60006122a8602a8361284d565b91506122b382612d0c565b604082019050919050565b60006122cb601e8361284d565b91506122d682612d5b565b602082019050919050565b60006122ee60298361284d565b91506122f982612d84565b604082019050919050565b600061231160298361284d565b915061231c82612dd3565b604082019050919050565b600061233460288361284d565b915061233f82612e22565b604082019050919050565b600061235760218361284d565b915061236282612e71565b604082019050919050565b6123768161291e565b82525050565b6123858161291e565b82525050565b60006020820190506123a060008301846120b8565b92915050565b600060a0820190506123bb60008301886120b8565b6123c860208301876120b8565b81810360408301526123da81866120c7565b905081810360608301526123ee81856120c7565b905081810360808301526124028184612134565b90509695505050505050565b600060a08201905061242360008301886120b8565b61243060208301876120b8565b61243d604083018661237c565b61244a606083018561237c565b818103608083015261245c8184612134565b90509695505050505050565b6000602082019050818103600083015261248281846120c7565b905092915050565b600060408201905081810360008301526124a481856120c7565b905081810360208301526124b881846120c7565b90509392505050565b60006020820190506124d66000830184612125565b92915050565b600060208201905081810360008301526124f6818461216d565b905092915050565b60006020820190508181036000830152612517816121a6565b9050919050565b60006020820190508181036000830152612537816121c9565b9050919050565b60006020820190508181036000830152612557816121ec565b9050919050565b600060208201905081810360008301526125778161220f565b9050919050565b6000602082019050818103600083015261259781612232565b9050919050565b600060208201905081810360008301526125b781612255565b9050919050565b600060208201905081810360008301526125d781612278565b9050919050565b600060208201905081810360008301526125f78161229b565b9050919050565b60006020820190508181036000830152612617816122be565b9050919050565b60006020820190508181036000830152612637816122e1565b9050919050565b6000602082019050818103600083015261265781612304565b9050919050565b6000602082019050818103600083015261267781612327565b9050919050565b600060208201905081810360008301526126978161234a565b9050919050565b60006020820190506126b3600083018461237c565b92915050565b60006040820190506126ce600083018561237c565b6126db602083018461237c565b9392505050565b60006126ec6126fd565b90506126f8828261299c565b919050565b6000604051905090565b600067ffffffffffffffff82111561272257612721612a74565b5b602082029050602081019050919050565b600067ffffffffffffffff82111561274e5761274d612a74565b5b602082029050602081019050919050565b600067ffffffffffffffff82111561277a57612779612a74565b5b602082029050602081019050919050565b600067ffffffffffffffff8211156127a6576127a5612a74565b5b6127af82612ac5565b9050602081019050919050565b600067ffffffffffffffff8211156127d7576127d6612a74565b5b6127e082612ac5565b9050602081019050919050565b6000819050602082019050919050565b600081519050919050565b600081519050919050565b600081519050919050565b6000602082019050919050565b600082825260208201905092915050565b600082825260208201905092915050565b600082825260208201905092915050565b60006128698261291e565b91506128748361291e565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff038211156128a9576128a8612a16565b5b828201905092915050565b60006128bf826128fe565b9050919050565b60008115159050919050565b60007fffffffff0000000000000000000000000000000000000000000000000000000082169050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b82818337600083830152505050565b60005b8381101561295557808201518184015260208101905061293a565b83811115612964576000848401525b50505050565b6000600282049050600182168061298257607f821691505b6020821081141561299657612995612a45565b5b50919050565b6129a582612ac5565b810181811067ffffffffffffffff821117156129c4576129c3612a74565b5b80604052505050565b60006129d88261291e565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415612a0b57612a0a612a16565b5b600182019050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600060033d1115612ac25760046000803e612abf600051612ad6565b90505b90565b6000601f19601f8301169050919050565b60008160e01c9050919050565b7f455243313135353a207472616e7366657220746f206e6f6e204552433131353560008201527f526563656976657220696d706c656d656e746572000000000000000000000000602082015250565b7f455243313135353a204552433131353552656365697665722072656a6563746560008201527f6420746f6b656e73000000000000000000000000000000000000000000000000602082015250565b7f5468652069647320616e6420616d6f756e747320617265206e6f74206d61746360008201527f6800000000000000000000000000000000000000000000000000000000000000602082015250565b7f455243313135353a2063616c6c6572206973206e6f74206f776e6572206e6f7260008201527f20617070726f7665640000000000000000000000000000000000000000000000602082015250565b7f455243313135353a2061646472657373207a65726f206973206e6f742061207660008201527f616c6964206f776e657200000000000000000000000000000000000000000000602082015250565b7f455243313135353a207472616e7366657220746f20746865207a65726f20616460008201527f6472657373000000000000000000000000000000000000000000000000000000602082015250565b7f455243313135353a207472616e736665722063616c6c6572206973206e6f742060008201527f6f776e6572206e6f7220617070726f7665640000000000000000000000000000602082015250565b7f455243313135353a20696e73756666696369656e742062616c616e636520666f60008201527f72207472616e7366657200000000000000000000000000000000000000000000602082015250565b7f5468652069647320616e64207572697320617265206e6f74206d617463680000600082015250565b7f455243313135353a2073657474696e6720617070726f76616c2073746174757360008201527f20666f722073656c660000000000000000000000000000000000000000000000602082015250565b7f455243313135353a206163636f756e747320616e6420696473206c656e67746860008201527f206d69736d617463680000000000000000000000000000000000000000000000602082015250565b7f455243313135353a2069647320616e6420616d6f756e7473206c656e6774682060008201527f6d69736d61746368000000000000000000000000000000000000000000000000602082015250565b7f455243313135353a206d696e7420746f20746865207a65726f2061646472657360008201527f7300000000000000000000000000000000000000000000000000000000000000602082015250565b600060443d1015612ed057612f53565b612ed86126fd565b60043d036004823e80513d602482011167ffffffffffffffff82111715612f00575050612f53565b808201805167ffffffffffffffff811115612f1e5750505050612f53565b80602083010160043d038501811115612f3b575050505050612f53565b612f4a8260200185018661299c565b82955050505050505b90565b612f5f816128b4565b8114612f6a57600080fd5b50565b612f76816128c6565b8114612f8157600080fd5b50565b612f8d816128d2565b8114612f9857600080fd5b50565b612fa48161291e565b8114612faf57600080fd5b5056fea2646970667358221220832b69e5ad295dcef45df51f723be8ec406eeae3928053f958e0568b9d4e087d64736f6c63430008010033"
	abi   = "[{\"inputs\": [],\"stateMutability\": \"nonpayable\",\"type\": \"constructor\"},{\"anonymous\": false,\"inputs\": [{\"indexed\": true,\"internalType\": \"address\",\"name\": \"account\",\"type\": \"address\"},{\"indexed\": true,\"internalType\": \"address\",\"name\": \"operator\",\"type\": \"address\"},{\"indexed\": false,\"internalType\": \"bool\",\"name\": \"approved\",\"type\": \"bool\"}],\"name\": \"ApprovalForAll\",\"type\": \"event\"},{\"anonymous\": false,\"inputs\": [{\"indexed\": true,\"internalType\": \"address\",\"name\": \"operator\",\"type\": \"address\"},{\"indexed\": true,\"internalType\": \"address\",\"name\": \"from\",\"type\": \"address\"},{\"indexed\": true,\"internalType\": \"address\",\"name\": \"to\",\"type\": \"address\"},{\"indexed\": false,\"internalType\": \"uint256[]\",\"name\": \"ids\",\"type\": \"uint256[]\"},{\"indexed\": false,\"internalType\": \"uint256[]\",\"name\": \"values\",\"type\": \"uint256[]\"}],\"name\": \"TransferBatch\",\"type\": \"event\"},{\"anonymous\": false,\"inputs\": [{\"indexed\": true,\"internalType\": \"address\",\"name\": \"operator\",\"type\": \"address\"},{\"indexed\": true,\"internalType\": \"address\",\"name\": \"from\",\"type\": \"address\"},{\"indexed\": true,\"internalType\": \"address\",\"name\": \"to\",\"type\": \"address\"},{\"indexed\": false,\"internalType\": \"uint256\",\"name\": \"id\",\"type\": \"uint256\"},{\"indexed\": false,\"internalType\": \"uint256\",\"name\": \"value\",\"type\": \"uint256\"}],\"name\": \"TransferSingle\",\"type\": \"event\"},{\"anonymous\": false,\"inputs\": [{\"indexed\": false,\"internalType\": \"string\",\"name\": \"value\",\"type\": \"string\"},{\"indexed\": true,\"internalType\": \"uint256\",\"name\": \"id\",\"type\": \"uint256\"}],\"name\": \"URI\",\"type\": \"event\"},{\"inputs\": [],\"name\": \"_owner\",\"outputs\": [{\"internalType\": \"address\",\"name\": \"\",\"type\": \"address\"}],\"stateMutability\": \"view\",\"type\": \"function\"},{\"inputs\": [{\"internalType\": \"address\",\"name\": \"account\",\"type\": \"address\"},{\"internalType\": \"uint256\",\"name\": \"id\",\"type\": \"uint256\"}],\"name\": \"balanceOf\",\"outputs\": [{\"internalType\": \"uint256\",\"name\": \"\",\"type\": \"uint256\"}],\"stateMutability\": \"view\",\"type\": \"function\"},{\"inputs\": [{\"internalType\": \"address[]\",\"name\": \"accounts\",\"type\": \"address[]\"},{\"internalType\": \"uint256[]\",\"name\": \"ids\",\"type\": \"uint256[]\"}],\"name\": \"balanceOfBatch\",\"outputs\": [{\"internalType\": \"uint256[]\",\"name\": \"\",\"type\": \"uint256[]\"}],\"stateMutability\": \"view\",\"type\": \"function\"},{\"inputs\": [{\"internalType\": \"address\",\"name\": \"account\",\"type\": \"address\"},{\"internalType\": \"address\",\"name\": \"operator\",\"type\": \"address\"}],\"name\": \"isApprovedForAll\",\"outputs\": [{\"internalType\": \"bool\",\"name\": \"\",\"type\": \"bool\"}],\"stateMutability\": \"view\",\"type\": \"function\"},{\"inputs\": [{\"internalType\": \"address\",\"name\": \"_to\",\"type\": \"address\"},{\"internalType\": \"uint256[]\",\"name\": \"ids\",\"type\": \"uint256[]\"},{\"internalType\": \"uint256[]\",\"name\": \"amounts\",\"type\": \"uint256[]\"},{\"internalType\": \"string[]\",\"name\": \"uris\",\"type\": \"string[]\"}],\"name\": \"mint\",\"outputs\": [],\"stateMutability\": \"nonpayable\",\"type\": \"function\"},{\"inputs\": [{\"internalType\": \"address\",\"name\": \"from\",\"type\": \"address\"},{\"internalType\": \"address\",\"name\": \"to\",\"type\": \"address\"},{\"internalType\": \"uint256[]\",\"name\": \"ids\",\"type\": \"uint256[]\"},{\"internalType\": \"uint256[]\",\"name\": \"amounts\",\"type\": \"uint256[]\"},{\"internalType\": \"bytes\",\"name\": \"data\",\"type\": \"bytes\"}],\"name\": \"safeBatchTransferFrom\",\"outputs\": [],\"stateMutability\": \"nonpayable\",\"type\": \"function\"},{\"inputs\": [{\"internalType\": \"address\",\"name\": \"from\",\"type\": \"address\"},{\"internalType\": \"address\",\"name\": \"to\",\"type\": \"address\"},{\"internalType\": \"uint256\",\"name\": \"id\",\"type\": \"uint256\"},{\"internalType\": \"uint256\",\"name\": \"amount\",\"type\": \"uint256\"},{\"internalType\": \"bytes\",\"name\": \"data\",\"type\": \"bytes\"}],\"name\": \"safeTransferFrom\",\"outputs\": [],\"stateMutability\": \"nonpayable\",\"type\": \"function\"},{\"inputs\": [{\"internalType\": \"address\",\"name\": \"operator\",\"type\": \"address\"},{\"internalType\": \"bool\",\"name\": \"approved\",\"type\": \"bool\"}],\"name\": \"setApprovalForAll\",\"outputs\": [],\"stateMutability\": \"nonpayable\",\"type\": \"function\"},{\"inputs\": [{\"internalType\": \"bytes4\",\"name\": \"interfaceId\",\"type\": \"bytes4\"}],\"name\": \"supportsInterface\",\"outputs\": [{\"internalType\": \"bool\",\"name\": \"\",\"type\": \"bool\"}],\"stateMutability\": \"view\",\"type\": \"function\"},{\"inputs\": [{\"internalType\": \"address\",\"name\": \"to\",\"type\": \"address\"},{\"internalType\": \"uint256\",\"name\": \"id\",\"type\": \"uint256\"},{\"internalType\": \"uint256\",\"name\": \"amount\",\"type\": \"uint256\"}],\"name\": \"transferArtNFT\",\"outputs\": [],\"stateMutability\": \"nonpayable\",\"type\": \"function\"},{\"inputs\": [{\"internalType\": \"uint256\",\"name\": \"_id\",\"type\": \"uint256\"}],\"name\": \"uri\",\"outputs\": [{\"internalType\": \"string\",\"name\": \"\",\"type\": \"string\"}],\"stateMutability\": \"view\",\"type\": \"function\"}]"
)

func TestEvm(t *testing.T) {
	jsonclient, err := client.NewJSONClient("", url)
	assert.Nil(t, err)

	// 部署合约
	code, err := types.FromHex(codes)
	assert.Nil(t, err)
	tx, err := CreateEvmContract(code, "", "evm-sdk-test", paraName)
	err = SignTx(tx, deployPrivateKey)
	assert.Nil(t, err)
	txhash := types.ToHexPrefix(tx.Hash())
	signTx := types.ToHexPrefix(types.Encode(tx))
	reply, err := jsonclient.SendTransaction(signTx)
	assert.Nil(t, err)
	assert.Equal(t, txhash, reply)
	fmt.Print("部署合约交易hash = ", txhash)
	time.Sleep(5 * time.Second)
	detail, err := jsonclient.QueryTransaction(txhash)
	assert.Nil(t, err)
	fmt.Println("; 执行结果 = ", detail.Receipt.Ty)

	// 计算合约地址
	contractAddress := crypto.GetExecAddress(deployAddress + strings.TrimPrefix(txhash, "0x"))
	fmt.Println("部署好的合约地址 = " + contractAddress)

	length := 2
	// tokenId数组
	ids := make([]int, length)
	// 同一个tokenid发行多少份
	amounts := make([]int, length)
	// 每一个tokenid对应的URI信息（一般对应存放图片的描述信息，图片内容的一个url）
	uris := make([]string, length)
	for i := 0; i < length; i++ {
		ids[i] = 10000 + i
		amounts[i] = 100
		// 例子为了简化处理，让所有ID都固定一个地址，
		uris[i] = "http://www.baidu.com"
	}
	idStr, _ := json.Marshal(ids)
	amountStr, _ := json.Marshal(amounts)
	uriStr, _ := json.Marshal(uris)

	// 调用合约
	param := fmt.Sprintf("mint(%s,%s,%s,%s)", useraAddress, idStr, amountStr, uriStr)
	initNFT, err := EncodeParameter(abi, param)
	assert.Nil(t, err)
	tx, err = CallEvmContract(initNFT, "", 0, contractAddress, paraName)
	// 构造交易组
	group, err := CreateNobalance(tx, useraPrivateKey, withholdPrivateKey, paraName)
	assert.Nil(t, err)
	txhash = types.ToHexPrefix(group.Tx().Hash())
	signTx = types.ToHexPrefix(types.Encode(group.Tx()))
	reply, err = jsonclient.SendTransaction(signTx)
	assert.Nil(t, err)
	assert.Equal(t, txhash, reply)
	fmt.Print("代扣交易hash = ", txhash)
	time.Sleep(5 * time.Second)
	detail, err = jsonclient.QueryTransaction(txhash)
	assert.Nil(t, err)
	fmt.Println("; 执行结果 = ", detail.Receipt.Ty)

	// 合约查询
	param = fmt.Sprintf("balanceOf(%s,%d)", useraAddress, ids[0])
	QueryContract(url, contractAddress, abi, param, contractAddress)
	param = fmt.Sprintf("uri(%d)", ids[1])
	QueryContract(url, contractAddress, abi, param, contractAddress)
}