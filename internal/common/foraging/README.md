## Foraging

### Deer Hunts 🦌

#### Utility for a single deer
The return for a single deer is modelled as a random variable that is implemented as `deer_return()`. It is effectively the combination of two other RVs:
- `D`: Bernoulli RV that represents the probability of catching a deer at all (binary). Usually `p` - i.e. P(`D`=1) = `p` - will be fairly close to 1 (fairly high chance of catching a deer if you invest the resources)
- `W`: A continuous RV that adds some variance to the return. This could be interpreted as the weight of the deer that is caught. W is exponentially distributed such that the prevalence of deer of certain size is inversely prop. to their size.

D and H are combined in the following expression: $H = D(1+W)$ (another RV which is the output of `deer_return()`). Notice how `W` is irrelevant if `D`=0 (weight of a deer does not matter if we don't catch it). Also, the mean return can be shown to be $E[U] = E[D(1+W)] = p(1+E[W]) = p(1+\frac{1}{\lambda})$ where $p$ is the Bernoulli parameter of `D` and $\lambda$ is the rate parameter of `W`

#### Dynamics of a collective hunt
A deer hunt is modelled as a stochastic process where the expected utility is proportional to input resources provided for the hunt. There is no limit to the number of teams `m` that can join a hunt, however there is a limit to the number of deer `n` that can be hunted in a single hunt. `n` should be chosen s.t. `n < max(m) = 6` so that beyond a certain number of hunters, the utility per team decreases as the same expected return could be achieved with less hunters. 

The utility distribution is modelled as being agnostic to the number of hunters; it only depends on the collective amount of resources invested in the hunt. Furthermore, assume the minimum resources to hunt a single deer is $\theta$. We use the (naive) assumption that hunting one deer is independent to another and let the expected return for a single deer be $E[U]=\mu$. Then, for a given real-valued collective resource input $x$ (across all hunt participants), the expected return is proposed as follows:

$$
U_{\theta}=\left\{\begin{array}{ll}
\mu & x \in[\theta ; 2 \theta) \\
2 \mu & x \in[2 \theta ; 2.75 \theta) \\
3 \mu & x \in[2.75\theta ; 3.25 \theta) \\
4 \mu & x \in[3.25\theta ; \infty)
\end{array}\right.
$$

where the maximum number of deer that can be hunted in a single hunt is 4. Notice that the expected return for $n$ deer is simply $n$ times the return of a single deer (i.i.d. assumption). *However*, the incremental input resources required to hunt $n$ deer *decreases* as $n$ increases. That is, it costs less (per deer) to hunt more deer. This is to incentivise **collaboration**. In this package, this dynamic is implemented in the `deerUtilityTier()` function that returns the number of deer that can be hunted (i.e. the utility tier) given a scalar collective resource input ($x$ from above). In this implementation, $\theta=1$, but it could be multiplied by an arbitrary multiplier to scale as desired.

#### ToDo

- fishing
- link deer population (governed by a predefined differential equation and historical consumption) with Bernoulli param `p` in deer return RV.