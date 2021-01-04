## Disasters :tornado:
*note*: GH doesn't support LaTeX in MD file previews. View this file locally in VSCode or similar to see the maths.

### Prevalence
Disasters are one of the key components to the the long term CRD. Disasters can be configured to occur with a **stochastic** or **deterministic** period $T_0$ by toggling the `StochasticPeriod` parameter in the `DisasterConfig`. These two cases are designed such that:
- In the *deterministic* case, a disaster occurs regularly with a period $T_0$ (i.e. it is guaranteed to occur every $T$ turns).
- In the *stochastic* case, the *expected* period $E[T]$ = $T_0$.

Note that in the stochastic case, the period is a *geometric* random variable as it models the number of turns before a disaster strikes (assuming individual disaster samples on each turn are independent). If a disaster occurs on a given turn with probability `p`, the expected value of this geometric RV, the period, = $E[T]$ = $1/p$. Since we want $E[T]$ = $T_0$, $p$ is implied when $T_0$ is given and so we only need to specify this period parameter to cover both cases.

### Severity and Location
In both the stochastic and deterministic cases, the **magnitude** and **location** of a disaster are sampled in the same fashion. The magnitude is exponentially distributed with scale parameter `ExponentialRate` in the `DisasterConig`. This was chosen to model a plausible real life scenario where smaller disasters are far more common than very serious ones. The xy co-ordinates of the *epicentre* (location of peak magnitude) of the disaster are sampled from a joint uniform distribution with bounds specified in the `DisasterConfig`. When a disaster strikes, the **effect** (damage) felt by a given island is inversely proportional to the square of its distance to the epicentre of the disaster.




