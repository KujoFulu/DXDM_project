library(shiny)
# ui.R

shinyUI(fluidPage(
  titlePanel("Wright-Fisher Simulation"),
  
  sidebarLayout(
    sidebarPanel(
      p("*For large parameters input it might take 20-30 seconds to display all plots. "),

      # All sliders
      sliderInput("popSize", "Population Size", min = 10, max = 500, value = 200),
      sliderInput("selCo", "Selection Coefficient", min = -0.2, max = 0.2, value = 0.02),
      sliderInput("freqStart", "Starting Allele Frequency", min = 0, max = 1, value = 0.5),
      sliderInput("numGen", "Number of Generations", min = 10, max =500, value = 100),
      sliderInput("numRuns", "Number of Runs", min = 1, max = 500, value =100),
      
      # Simulation button
      actionButton("runButton", "Run Simulation"),

      # Saving button
      actionButton("savePlotsButton", "Save All Plots"),
      
      # Simulation without Plotting button
      actionButton("runSimulationWithoutPlotting", "Simulation without Plotting")
    ),
    
    mainPanel(
      withMathJax(),
      h4("This is a simulator that can simulate the haploid gene drift with Wright-Fisher model. "),
      p("In this sigle site Wright-Fisher model we have panmictic population size N over t generations for a single loci with two alleles a and A."),
      p("n is the start allele number of A which also indicate the product of start allele frequency of A and population size"),
      p("The number of alleles 𝑛′ in each generation is sampled independently from a binomial distribution. The selection coefficient 𝑠 can affect the selection of alleles. When 𝑠 is positive, allele A is more advantageous and easier to retain, while a negative 𝑠 means that the allele is harmful and difficult to pass on to the next generation."),
      p("This is the number of allele A's function in each generation"),
      p(withMathJax("$$n' \\sim binomial(\\frac{n(1 + s)}{n(1 + s) + N - n}, N)$$")),
      p("When the proportion of allele A is 100% we call it fix, on the contrary if the proportion is 0% we call it loss."),
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
