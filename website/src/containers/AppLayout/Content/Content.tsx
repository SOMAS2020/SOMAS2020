import React from "react"
import { Route, Switch } from "react-router-dom"
import Home from "../../../components/Home/Home"
import * as Visualisations from "../../../components/Visualisations"
import RawOutput from "../../../components/RawOutput/RawOutput"

const Content = () => {
  return (
    <div>
      <Switch>
        <Route path="/rawoutput" exact component={RawOutput}/>
        <Route path="/resources" component={Visualisations.Resources}/>
        <Route path="/roles"/>
        <Route path="/IIGO"/>
        <Route path="/IITO"/>
        <Route path="/IIFO"/>
        <Route path="/" component={Home}/>
      </Switch>
    </div>
  )
}

export default Content