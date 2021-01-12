import React from 'react'
import { withStyles, makeStyles } from '@material-ui/core/styles'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableContainer from '@material-ui/core/TableContainer'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'
import Paper from '@material-ui/core/Paper'

const useStyles = makeStyles({
  table: {
    minWidth: 650,
  },
})

function createData(name: string, AgentDescription: string) {
  return { name, AgentDescription }
}

const Agent1Desc =
  "For team 1's agent, our strategy for foraging is to always forage for deer with an amount inversely proportional to the sum of last turn's foraging. When the deer population is low, our agent will contribute a small amount of resources to forage for fish to allow for the deer population to regenerate. Our agent also tries to follow the rules by paying tax and requesting allocations from the commonpool when it is in a critical state (meaning it has no resources left). For gifts, our agent requests gifts only when the agent is in a critical state. On the other hand, when not in a  critical state, our agent will give gifts depending on our agent's opinion of an island. Our agent's opinion of an island will increase if the island offers us gifts."

const Agent2Desc =
  "Our agent will counteract the strategies of the other agents in order to keep the common pool stable. For example, if the other agents are free riders on average, then our agent will become an altruist for a limited number of turns in order to try and help save the common pool and the game. We apply a strategy of 'herd mentality' in terms of resource allocation, by only contributing factors of the average contribution to the pool depending on situations and our own internal state."

const Agent3Desc =
  "Team 3's agent has an emphasis on being highly parameterised. The behaviour of each function depends on the value of the parameters. This allows the agent to work in a broader scenario and also allow multiple behaviours of agents by simply changing the parameters. Our agent also includes geometric analysis of the matrix based rules, allowing us to more accurately make quantitative judgements on the effects of newly proposed rules on our agents performance. We use the same geometric analysis to propose modifications to rules to better suit our current circumstances. This system is combined with our parameterized compliance approach to investigate elements of game theory and are useful in the study of exploitation of the IIGO system."

const Agent4Desc =
  "We have implemented a truthful agent, which is only be able to deviate from following rules based on the personality parameters it might get. It is a 'vanilla' agent implementation, which does not do anything based on randomness, but rather rules."

const Agent5Desc =
  "Team 5's agent is driven by historical data to form opinions on other agents. This knowledge formation is used to evaluate future steps of action that include foraging decisions, gifting decisions, common pool contributions and role management, as well as disaster forecasting. This enables our agent to engage in more complex yet more stochastic decision making. Our agent can be characterised as fair and generous, and it has been designed to support the rest of the agent community when necessary. The latter is particularly true when there is open communication of data between islands. However, our agent can be biased at times in order to ensure our island's survival when resources are scarce."

const Agent6Desc =
  'Our decision-making is based on economic and relationship aspects. Depending on the resources we own, we act in different personalities in the game of investment. Likewise, we trust our friends and help each other in the game of sustainability. We evaluate the collective precautionary measure by comparing taxation with our prediction of disaster, and always advocate for a surplus of common-pool contribution. When we are elected for roles, we take the responsibility to maintain fairness and justice, with the premise that our friends benefit from all the moves.'

const rows = [
  createData('Team 1', Agent1Desc),
  createData('Team 2', Agent2Desc),
  createData('Team 3', Agent3Desc),
  createData('Team 4', Agent4Desc),
  createData('Team 5', Agent5Desc),
  createData('Team 6', Agent6Desc),
]

const AgentTable = () => {
  const classes = useStyles()
  return (
    <TableContainer component={Paper}>
      <Table
        className={classes.table}
        size="small"
        aria-label="customized table"
      >
        <colgroup>
          <col style={{ width: '10%' }} />
          <col style={{ width: '90%' }} />
        </colgroup>
        <TableHead>
          <TableRow style={{ backgroundColor: 'lightblue', color: 'white' }}>
            <TableCell>Team Number</TableCell>
            <TableCell align="left">Agent Description</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {rows.map((row) => (
            <TableRow
              key={row.name}
              style={{ backgroundColor: 'white', color: 'white' }}
            >
              <TableCell component="th" scope="row">
                {row.name}
              </TableCell>
              <TableCell align="left">{row.AgentDescription}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  )
}

export default AgentTable
