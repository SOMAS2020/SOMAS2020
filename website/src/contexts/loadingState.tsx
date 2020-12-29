import React from 'react'

const { createContext, useReducer, useContext } = React;

export type LoadingState = {
  loading: boolean
  loadingText?: string | undefined
}

export type LoadingStateDispatchType = (values: LoadingState) => void;

export const initialLoadingState: LoadingState =
  { loading: false, loadingText: undefined }

const LoadingStateContext = createContext(initialLoadingState)
const DispatchLoadingStateContext = createContext({} as LoadingStateDispatchType)


export const LoadingStateProvider = ({ children }: { children: any }) => {
  const [state, dispatch] = useReducer(
    (_: LoadingState, newValue: LoadingState) => {
      return newValue
    },
    initialLoadingState
  )
  return (
    <LoadingStateContext.Provider value={state}>
      <DispatchLoadingStateContext.Provider value={dispatch}>
        {children}
      </DispatchLoadingStateContext.Provider>
    </LoadingStateContext.Provider>
  )
}

export const useLoadingState = (): [LoadingState, LoadingStateDispatchType] => {
  return [
    useContext(LoadingStateContext),
    useContext(DispatchLoadingStateContext)
  ]
}