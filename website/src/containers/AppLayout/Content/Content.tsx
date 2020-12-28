import React from "react"
import { Route, Switch } from "react-router-dom"
import Home from "../../../components/Home/Home"
import CIOutput from "../../../components/CIOutput/CIOutput"
import { cioutput } from "../../../consts/paths"

const Content = () => {
  return (
    <div>
      <Switch>
        <Route path={cioutput} exact component={CIOutput}/>
        <Route component={Home}/>
      </Switch>
    </div>
  )
}

export default Content