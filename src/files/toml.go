package files
import(
	"github.com/Kdag-K/kdag-hub/src/common"
	"github.com/pelletier/go-toml"
)


// LoadToml loads a Toml file and returns a tree object.
func LoadToml(tomlFile string) (*toml.Tree, error) {
	config, err := toml.LoadFile(tomlFile)

	if err != nil {
		common.ErrorMessage("Error loading toml file: ", tomlFile)
		return nil, err
	}

	return config, nil
}
