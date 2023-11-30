library(shiny)
# ui.R

shinyUI(fluidPage(
  titlePanel("Wright-Fisher Simulation"),
  
  sidebarLayout(
    sidebarPanel(
      p("This is a simulator that can simulate the gene drift with Wright-Fisher model.For display plots quickly recommed generation<200."),
      sliderInput("popSize", "Population Size", min = 10, max = 500, value = 200),
      sliderInput("selCo", "Selection Coefficient", min = -0.2, max = 0.2, value = 0.02),
      sliderInput("freqStart", "Starting Allele Frequency", min = 0, max = 1, value = 0.5),
      sliderInput("numGen", "Number of Generations", min = 10, max =500, value = 100),
      sliderInput("numRuns", "Number of Runs", min = 1, max = 500, value =100),
      actionButton("runButton", "Run Simulation"),

      # Saving button
      actionButton("savePlotsButton", "Save All Plots"),
      
      # Simulation without Plotting button
      actionButton("runSimulationWithoutPlotting", "Simulation without Plotting")
    ),
    
    mainPanel(
      tabsetPanel(
        tabPanel("Data", 
                 tableOutput("simulationData")),
        tabPanel("Average Allele Frequency", 
                 plotOutput("average_allele_freq_plot")),
        tabPanel("Genotype Frequency", 
                 plotOutput("genotype_frequency_plot")),
        tabPanel("Total Allele Copies Histogram", 
                 plotOutput("totalAlleleCopiesHistogram")),
        tabPanel("Fixation and Loss Bar Plots", 
                 plotOutput("fix_loss_bar_plot")),
        tabPanel("Max Allele Histogram", 
                 plotOutput("max_allele_histogram")),
        tabPanel("Heat Plot", 
                 plotOutput("combined_plot")),
        tabPanel("Allele Fixed and Lossed Ratio", 
                 plotOutput("lineplot_flratio")),
        tabPanel("Cumulative Line Plot Animation", 
                 plotlyOutput("fig_genotype"),
                 hr()),
        tabPanel("Allele A Frequency in 5 Runs", 
                 plotlyOutput("fig_5runs"))
      )
    )
  )
))
