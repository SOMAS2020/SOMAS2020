import React from 'react'
import { Button } from 'react-bootstrap'
import { useHistory } from 'react-router-dom'
import { useLoadingState } from '../../contexts/loadingState'
import { OutputJSONType } from '../../consts/types'
import { storeLocalVisOutput } from '../Visualisations/utils'
import { visualisations } from '../../consts/paths'

const VisualiseButton = (props: { output: OutputJSONType }) => {
  const [, setLoading] = useLoadingState()
  const history = useHistory()

  const handleClick = async () => {
    setLoading({ loading: true, loadingText: `I can show you the world` })

    // sleep for 100ms for stability
    await new Promise((r) => setTimeout(r, 100))

    const { output } = props
    await storeLocalVisOutput(output)
    history.push(visualisations)
    // don't need to unload as after pushing to vis, there is another loading
    // instance
  }

  return (
    <>
      <Button variant="success" size="lg" onClick={handleClick}>
        Visualise
      </Button>
    </>
  )
}

export default VisualiseButton
