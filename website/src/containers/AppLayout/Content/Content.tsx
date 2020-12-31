import React from "react";
import { Route, Switch } from "react-router-dom";
import {
  cioutput,
  newrun,
  gamevisualisation,
  iigovisualisation,
  iitovisualisation,
  iifovisualisation,
  rolesvisualisation,
  resourcesvisualisation,
} from "../../../consts/paths";

import Home from "../../../components/Home/Home";
import CIOutput from "../../../components/CIOutput/CIOutput";
import NewRun from "../../../components/NewRun/NewRun";
import * as Visualisations from "../../../components/Visualisations/Visualisations";

const Content = () => {
  return (
    <div>
      <Switch>
        <Route path={cioutput} exact component={CIOutput} />
        <Route path={newrun} exact component={NewRun} />
        <Route path={gamevisualisation} exact component={Visualisations.Game} />
        <Route path={iigovisualisation} exact component={Visualisations.IIGO} />
        <Route path={iitovisualisation} exact component={Visualisations.IITO} />
        <Route path={iifovisualisation} exact component={Visualisations.IIFO} />
        <Route
          path={rolesvisualisation}
          exact
          component={Visualisations.Roles}
        />
        <Route
          path={resourcesvisualisation}
          exact
          component={Visualisations.Resources}
        />
        <Route component={Home} />
      </Switch>
    </div>
  );
};

export default Content;
