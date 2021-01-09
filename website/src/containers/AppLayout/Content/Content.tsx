import React from 'react'
import { Route, Switch } from 'react-router-dom'
import { about, cioutput, newrun, visualisations } from '../../../consts/paths'

import Home from '../../../components/Home/Home'
import About from '../../../components/About/About'
import CIOutput from '../../../components/CIOutput/CIOutput'
import NewRun from '../../../components/NewRun/NewRun'
import Visualisations from '../../../components/Visualisations/Visualisations'
import FourOhFour from '../../../components/FourOhFour/FourOhFour'

const Content = () => {
  return (
    <div>
      <Switch>
        <Route path={about} exact component={About} />
        <Route path={cioutput} exact component={CIOutput} />
        <Route path={newrun} exact component={NewRun} />
        <Route path={visualisations} component={Visualisations} />
        <Route path="/" exact component={Home} />
        <Route component={FourOhFour} />
      </Switch>
    </div>
  )
}

export default Content
