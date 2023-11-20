# DXDM_project
 02601_final_project

Wright-Fisher model simulation is a program that uses Go language to simulate and R language for visualization.

Before starting the simulation, some required packages must be added to the Go path: gonum.org/v1/gonum/stat/distuv
To install this package, you can use this code:
go install gonum.org/v1/gonum/stat/distuv@latest

The other packages in both Go and R can be automatically installed if they do not exist.
If the package load fails, please try to manually add packages. Here are all the packages required:
Go: 	"encoding/csv", "fmt",		"log", "time", "os",	"path/filepath", "gonum.org/v1/gonum/stat/distuv"
R: "ggplot2", "ggpubr", "dplyr", "viridis"
