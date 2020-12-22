import React from "react"
import { Route, Switch } from "react-router-dom"
import Home from "../../../components/Home/Home"
import RawOutput from "../../../components/RawOutput/RawOutput"

const Content = () => {
  return (
    <div>
      <Switch>
        <Route path="/rawoutput" exact component={RawOutput}/>
        <Route component={Home}/>
      </Switch>
    </div>
  )
}

export default Content