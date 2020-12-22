## Numerical Simulation

### ODEs

Differential equations are extremely useful for modelling evolving quantities such as animal populations that change in response to some stimulus. This package provides a basic ODE (ordinary DE) solver interface using the 4th order Runge-Kutta method (**RK4**). Read more on this [here](https://en.wikipedia.org/wiki/Runge–Kutta_methods#The_Runge–Kutta_method).

#### Practical application: modelling the deer population
As implemented in the `foraging` package, a population can be represented by a simple DE of the following form: 
$$\frac{dP}{dt} = k(N-P(t)) $$ 
where $k$ is the 'growth coefficient', $N$ is the carrying capacity (max population size) of the environment and $P$ is the instantaneous population size. This is similar to a classic logistic population model; the rate of growth of the population slows as it reaches carrying capacity.