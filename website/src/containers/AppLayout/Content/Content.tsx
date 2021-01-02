import React from "react";
import { Route, Switch } from "react-router-dom";
import {
  cioutput,
  newrun,
  visualisations,
} from "../../../consts/paths";

import Home from "../../../components/Home/Home";
import CIOutput from "../../../components/CIOutput/CIOutput";
import NewRun from "../../../components/NewRun/NewRun";
import Visualisations from "../../../components/Visualisations/Visualisations";

const Content = () => {
  return (
    <div>
      <Switch>
        <Route path={cioutput} exact component={CIOutput} />
        <Route path={newrun} exact component={NewRun} />
        <Route path={visualisations} component={Visualisations} />
        <Route component={Home} />
      </Switch>
    </div>
  );
};

export default Content;
