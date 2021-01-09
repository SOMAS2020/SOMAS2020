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

// const StyledTableCell = withStyles((theme) => ({
//   head: {
//     backgroundColor: theme.palette.common.blue,
//     color: theme.palette.common.black,
//   },
//   body: {
//     fontSize: 14,
//   },
// }))(TableCell)

// const StyledTableRow = withStyles((theme) => ({
//   root: {
//     '&:nth-of-type(odd)': {
//       backgroundColor: theme.palette.action.hover,
//     },
//   },
// }))(TableRow)

function createData(name: string, AgentDescription: string) {
  return { name, AgentDescription }
}

const Agent2Desc =
  'Our agent uses an approach based on evolutionary economic theory. We identify the others islands on average as fair sharers, free-riders or altruists and then decide which of these three roles we should play as to ensure our best chance of survival. We also consider the current and past states of the game to choose our actions. '

const Agent3Desc =
  'Team 3 agent has an emphasis on being highly parameterised. The behaviour of each function depends on the value of the parameters. This allows to agent to work in a broader scenario and also allow multiple behaviours of agents by simply changing the parameters.'

const rows = [
  createData('Team 1', 'Description of agent'),
  createData('Team 2', Agent2Desc),
  createData('Team 3', Agent3Desc),
  createData('Team 4', 'Description of agent'),
  createData('Team 5', 'Description of agent'),
  createData('Team 6', 'Description of agent'),
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
          <col style={{ width: '15%' }} />
          <col style={{ width: '85%' }} />
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
