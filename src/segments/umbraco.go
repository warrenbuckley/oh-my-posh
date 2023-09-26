package segments

import (
	"encoding/xml"
	"path/filepath"
	"strings"

	"github.com/jandedobbeleer/oh-my-posh/src/platform"
	"github.com/jandedobbeleer/oh-my-posh/src/properties"
)

const (
	umbracoFolderName = "umbraco"
)

type CSProj struct {
	PackageReferences []struct {
		Name    string `xml:"Include,attr"`
		Version string `xml:"Version,attr"`
	} `xml:"ItemGroup>PackageReference"`
}

// type CSProj struct {
// 	PackageReferences []PackageReference `xml:"ItemGroup>PackageReference"`
// }

// type PackageReference struct {
// 	Name    string `xml:"Include,attr"`
// 	Version string `xml:"Version,attr"`
// }

type Umbraco struct {
	props properties.Properties
	env   platform.Environment

	FoundUmbraco    bool
	IsModernUmbraco bool
	IsLegacyUmbraco bool
	Version         string
}

// Create a struct to use with XML Unmrashal for CSProj files to find PackageReferences items

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

			// TODO use XML unmarshal on contents
			csProjPackages := CSProj{}
			err := xml.Unmarshal([]byte(contents), &csProjPackages)

			if err != nil {
				// Log an error
			}

			// Loop over all the package references
			for _, packageReference := range csProjPackages.PackageReferences {
				if strings.ToLower(packageReference.Name) == strings.ToLower("umbraco.cms") {
					u.IsModernUmbraco = true
					u.FoundUmbraco = true

					u.Version = packageReference.Version
					u.env.Debug("UMBRACO: Found Umbraco.Cms in " + file)
					return true
				}
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
}

func (u *Umbraco) Template() string {
	return "UMBRACO !"
}

func (u *Umbraco) Init(props properties.Properties, env platform.Environment) {
	u.props = props
	u.env = env
}

// func (u *UserConfig) getEndpoint(name string) *EndpointConfig {
// 	endpoint, exists := u.Endpoints[name]

// 	if exists {
// 		return &endpoint
// 	}

// 	return nil
// }
