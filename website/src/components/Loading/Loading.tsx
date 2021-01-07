import React from 'react'
import { Modal } from 'react-bootstrap'
import logo from '../../assets/logo/logo512.png'
import { useLoadingState } from '../../contexts/loadingState'

import styles from './Loading.module.css'

const Loading = () => {
  const [loadingState] = useLoadingState()
  const { loading, loadingText } = loadingState

  return (
    <Modal
      size="lg"
      aria-labelledby="contained-modal-title-vcenter"
      centered
      show={loading}
      onHide={() => {}}
      style={{ textAlign: `center`, margin: `auto` }}
    >
      <Modal.Header>
        <Modal.Title id="contained-modal-title-vcenter">Loading...</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <img src={logo} className={styles.appLogo} alt="logo" />
        <p style={{ padding: 10, fontSize: '1.2rem' }}>
          {loadingText || `Tomato sauce was sold in the 1800's as medicine.`}
        </p>
      </Modal.Body>
    </Modal>
  )
}

export default Loading
