package segments

import (
	"path/filepath"
	"strings"

	"github.com/jandedobbeleer/oh-my-posh/src/platform"
	"github.com/jandedobbeleer/oh-my-posh/src/properties"
)

const (
	umbracoFolderName = "umbraco"
	// userFileName       = "user.json"
	// defaultEnpointName = "default"
)

type Umbraco struct {
	props properties.Properties
	env   platform.Environment

	FoundUmbraco    bool
	IsModernUmbraco bool
	IsLegacyUmbraco bool

	// Was sitecore stuff
	EndpointName string
	CmHost       string
}

// type EndpointConfig struct {
// 	Host string `json:"host"`
// }

// type UserConfig struct {
// 	DefaultEndpoint string                    `json:"defaultEndpoint"`
// 	Endpoints       map[string]EndpointConfig `json:"endpoints"`
// }

func (u *Umbraco) Enabled() bool {
	u.env.Debug("UMBRACO: Checking if we enable segment")

	// If the cwd does not contain a folder called 'umbraco'
	// Then get out of here...
	if !u.env.HasFolder(umbracoFolderName) {
		return false
	}

	// Check if we have a .csproj OR a web.config in the CWD
	if !u.env.HasFiles("*.csproj") && !u.env.HasFiles("web.config") {
		u.env.Debug("UMBRACO: NO CSProj or web.config found")
		return false
	}

	// Modern .NET Core based Umbraco
	if u.env.HasFiles("*.csproj") {
		u.env.Debug("UMBRACO: Found one or more .csproj files")

		// Open file contents and look for Umbraco.Cms
		// But there is no guranatee that the user my have commented it out
		// <!-- <PackageReference Include="Umbraco.Cms" Version="12.1.2"/> -->
		// <PackageReference Include="Umbraco.Cms" Version="12.1.2"/>

		// Find all .csproj files
		// searchDir := "."
		searchPattern := "*.csproj"

		// Get a list of all files that match the search pattern
		files, err := filepath.Glob(searchPattern)

		if err != nil {
			u.env.Debug("UMBRACO: Error while searching for .csproj files")
			u.env.Debug(err.Error())
			return false
		}

		// Loop over all the files that have a .csproj extension
		for _, file := range files {
			u.env.Debug("UMBRACO: Trying to open file at " + file)

			// Read the file contents of the csproj file
			contents := u.env.FileContent(file)

			if strings.Contains(contents, "PackageReference Include=\"Umbraco.Cms\"") {
				u.IsModernUmbraco = true
				u.FoundUmbraco = true
				u.env.Debug("UMBRACO: Found Umbraco.Cms in " + file)
				return true
			} else {
				u.env.Debug("UMBRACO: Found file but not the contents")
			}
		}
	} else {
		u.env.Debug("UMBRACO: SAD face")
		u.FoundUmbraco = false
		return false
	}

	// Got here then we should have returned true by now...
	u.FoundUmbraco = false
	return false

	// if !u.env.HasFiles(sitecoreFileName) || !u.env.HasFiles(path.Join(sitecoreFolderName, userFileName)) {
	// 	return false
	// }

	// var userConfig, err = getUserConfig(u)

	// if err != nil {
	// 	return false
	// }

	// u.EndpointName = userConfig.getDefaultEndpoint()

	// displayDefault := u.props.GetBool(properties.DisplayDefault, true)

	// if !displayDefault && u.EndpointName == defaultEnpointName {
	// 	return false
	// }

	// if endpoint := userConfig.getEndpoint(u.EndpointName); endpoint != nil && len(endpoint.Host) > 0 {
	// 	u.CmHost = endpoint.Host
	// }
}

func (u *Umbraco) Template() string {
	return "UMBRACO !"
}

func (u *Umbraco) Init(props properties.Properties, env platform.Environment) {
	u.props = props
	u.env = env
}

// func getUserConfig(s *Sitecore) (*UserConfig, error) {
// 	userJSON := s.env.FileContent(path.Join(sitecoreFolderName, userFileName))
// 	var userConfig UserConfig

// 	if err := json.Unmarshal([]byte(userJSON), &userConfig); err != nil {
// 		return nil, err
// 	}

// 	return &userConfig, nil
// }

// func (u *UserConfig) getDefaultEndpoint() string {
// 	if len(u.DefaultEndpoint) > 0 {
// 		return u.DefaultEndpoint
// 	}

// 	return defaultEnpointName
// }

// func (u *UserConfig) getEndpoint(name string) *EndpointConfig {
// 	endpoint, exists := u.Endpoints[name]

// 	if exists {
// 		return &endpoint
// 	}

// 	return nil
// }
