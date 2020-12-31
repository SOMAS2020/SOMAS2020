import React from "react"
import Home from "../../../components/Home/Home"
import { Route, Switch } from "react-router-dom"
import {
  cioutput,
  newrun,
  gamevisualisation,
  iigovisualisation,
  iitovisualisation,
  iifovisualisation,
  rolesvisualisation,
  resourcesvisualisation
} from "../../../consts/paths"

import CIOutput from "../../../components/CIOutput/CIOutput"
import NewRun from "../../../components/NewRun/NewRun"
import GameVisualisation from "../../../components/visualisation/Game/Game"
import IIGOVisualisation from "../../../components/visualisation/IIGO/IIGO"
import IITOVisualisation from "../../../components/visualisation/IITO/IITO"
import IIFOVisualisation from "../../../components/visualisation/IIFO/IIFO"
import RolesVisualisation from "../../../components/visualisation/Roles/Roles"
import ResourcesVisualisation from "../../../components/visualisation/Resources/Resources"

const Content = () => {
  return (
    <div>
      <Switch>
        <Route path={cioutput} exact component={CIOutput} />
        <Route path={newrun} exact component={NewRun} />
        <Route path={gamevisualisation} exact component={GameVisualisation} />
        <Route path={iigovisualisation} exact component={IIGOVisualisation} />
        <Route path={iitovisualisation} exact component={IITOVisualisation} />
        <Route path={iifovisualisation} exact component={IIFOVisualisation} />
        <Route path={rolesvisualisation} exact component={RolesVisualisation} />
        <Route path={resourcesvisualisation} exact component={ResourcesVisualisation} />
        <Route component={Home} />
      </Switch>
    </div>
  )
}

export default Content;