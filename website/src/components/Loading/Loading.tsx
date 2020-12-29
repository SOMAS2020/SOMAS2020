import React from 'react'
import { Modal } from 'react-bootstrap'

// import loadingSVG from './Loading.svg'
import { useLoadingState } from '../../contexts/loadingState'

const Loading = () => {
  const [loadingState] = useLoadingState()
  const { loading, loadingText } = loadingState

  return (
    <Modal
      size="lg"
      aria-labelledby="contained-modal-title-vcenter"
      centered
      show={loading}
      style={{ textAlign: `center`, margin: `auto` }}
    >
      <Modal.Header>
        <Modal.Title id="contained-modal-title-vcenter">
          Loading...
       </Modal.Title>
      </Modal.Header>
      <Modal.Body>
        {/* TODO:- Fix Freezing issue--might need web worker? */}
        {/* <img src={loadingSVG} alt="loading" className="loader" /> */}
        <p>
          <div style={{ padding: 10, fontSize: '1.2rem' }}>{loadingText ? loadingText : `Tomato sauce was sold in the 1800's as medicine.`}</div>
        </p>
      </Modal.Body>
    </Modal>
  )
}

export default Loading