library(plotly)
library(shiny)
library(shinyMatrix)  # for matrixInput

# Define UI
ui <- fluidPage(
  titlePanel("Dynamic Multi-Species Predator-Prey Model"),
  
  # Sidebar layout
  sidebarLayout(
    sidebarPanel(
      
      # Add a title
      h3("Examples from Scientific Literature"),
      # Add a button for running simulation with default values from paper
      actionButton("original", "Simple Three-Species"),
      actionButton("chaotic", "Chaotic Dynamics"),
      actionButton("extinctionOne", "Extinction of One Species"),
      actionButton("extinctionTwo", "Extinction of Two Species"),
      actionButton("limit", "Limit Cycle"),
      actionButton("stable", "Stable Equilibrium"),
      actionButton("lynx", "Classic Lynx-Hare Data"),
      
      # Add a title
      h3("Customize Your Own Ecosystem!"),
      # Input for number of species
      numericInput("numSpecies", 
                   "Number of Species:", 
                   min = 2, value = 1),
      
      # Dynamic UI for sliders and matrix
      uiOutput("dynamicUI"),
      
      # Action buttons
      actionButton("runSimulation", "Run Simulation"),
      actionButton("savePlots", "Save All Plots"),
      
    ),
    
    
    # Main panel for displaying outputs
    mainPanel(
      tabsetPanel(
        # display the first 50 rows of the simulation data csv
        tabPanel("Data", 
                 tableOutput("simulationData")),
        tabPanel("Animated Gif", 
                 plotOutput("animated_gif")),
        tabPanel("Population Dynamics Over Time", 
                 plotOutput("time_plot")),
        # only for the first two species
        tabPanel("Phase Plane Plot", 
                 plotOutput("phase_plane_plot")),
        tabPanel("Vector Field Plot", 
                 plotOutput("vector_field_plot")),
        tabPanel("Contour Plot", 
                 plotOutput("contour_plot"))
      )
    )
  )
)


