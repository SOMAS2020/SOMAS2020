import React from 'react'
// UI imports
import { makeStyles } from '@material-ui/core/styles'
import Box from '@material-ui/core/Box'
import Collapse from '@material-ui/core/Collapse'
import IconButton from '@material-ui/core/IconButton'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableContainer from '@material-ui/core/TableContainer'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'
import Typography from '@material-ui/core/Typography'
import Paper from '@material-ui/core/Paper'
import KeyboardArrowDownIcon from '@material-ui/icons/KeyboardArrowDown'
import KeyboardArrowUpIcon from '@material-ui/icons/KeyboardArrowUp'
import List from '@material-ui/core/List'
import ListItemText from '@material-ui/core/ListItemText'

import styles from './IIGO.module.css'
import { OutputJSONType } from '../../../consts/types'
import { RuleType, processRulesData } from './Util/ProcessedData'

const useRowStyles = makeStyles({
  root: {
    '& > *': {
      borderBottom: 'unset',
    },
  },
})

function Row(props: { row: RuleType }) {
  const { row } = props
  const [open, setOpen] = React.useState(false)
  const classes = useRowStyles()

  return (
    <>
      <TableRow className={classes.root}>
        <TableCell>
          <IconButton
            aria-label="expand row"
            size="small"
            onClick={() => setOpen(!open)}
          >
            {open ? <KeyboardArrowUpIcon /> : <KeyboardArrowDownIcon />}
          </IconButton>
        </TableCell>
        <TableCell component="th" scope="row">
          {row.ruleName}
        </TableCell>
        <TableCell align="right">{row.mutable ? 'True' : 'False'}</TableCell>
        <TableCell align="right">{row.linked ? 'True' : 'False'}</TableCell>
      </TableRow>
      <TableRow>
        <TableCell style={{ paddingBottom: 0, paddingTop: 0 }} colSpan={6}>
          <Collapse in={open} timeout="auto" unmountOnExit>
            <Box margin={1}>
              <Typography variant="h6" gutterBottom component="div">
                Required Variables
                <List component="nav">
                  {row.variables.map((variable) => (
                    <ListItemText key={variable} inset secondary={variable} />
                  ))}
                </List>
              </Typography>
              <Typography variant="h6" gutterBottom component="div">
                History
              </Typography>
              <Table size="small" aria-label="purchases">
                <TableHead>
                  <TableRow>
                    <TableCell>Season</TableCell>
                    <TableCell align="right">Turn</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {row.history.map((historyRow) => (
                    <TableRow key={historyRow.season * 365 + historyRow.turn}>
                      <TableCell component="th" scope="row">
                        {historyRow.season}
                      </TableCell>
                      <TableCell align="right">{historyRow.turn}</TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </Box>
          </Collapse>
        </TableCell>
      </TableRow>
    </>
  )
}

const IIGO = (props: { output: OutputJSONType }) => {
  const rows: RuleType[] = processRulesData(props.output)

  return (
    <div className={styles.root}>
      <p className={styles.text}>Rules</p>
      <TableContainer
        style={{ maxHeight: 500, overflow: 'auto' }}
        component={Paper}
      >
        <Table aria-label="collapsible table">
          <TableHead>
            <TableRow>
              <TableCell />
              <TableCell>Rule Name</TableCell>
              <TableCell align="right">Mutable</TableCell>
              <TableCell align="right">Linked</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {rows.map((row) => (
              <Row key={row.ruleName} row={row} />
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </div>
  )
}

export default IIGO
