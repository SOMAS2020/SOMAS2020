import React from "react"
import { Route, Switch } from "react-router-dom"
import Home from "../../../components/Home/Home"
import CIOutput from "../../../components/CIOutput/CIOutput"
import NewRun from "../../../components/NewRun/NewRun"
import { cioutput, newrun } from "../../../consts/paths"

const Content = () => {
  return (
    <div>
      <Switch>
        <Route path={cioutput} exact component={CIOutput}/>
        <Route path={newrun} exact component={NewRun}/>
        <Route component={Home}/>
      </Switch>
    </div>
  )
}

export default Content