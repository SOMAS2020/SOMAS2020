import React from "react"
import { Route, Switch } from "react-router-dom"
import Home from "../../../components/Home/Home"
import { cioutput } from "../../../consts/paths"
import CIOutput from "../../../components/CIOutput/CIOutput"

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