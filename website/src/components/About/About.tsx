import React from 'react'
import styles from './About.module.css'
import AgentTable from './AgentTable'

const About = () => {
  return (
    <div style={{ paddingTop: 50 }}>
      <h1>Self-Organising Multi-Agent Systems 2020</h1>
      <div className={styles.root}>
        <p>
          This project explores different areas of self-organising multi-agent
          systems (SOMAS) through a game of survival. In this game, there are
          multiple islands that must work together to survive disasters by
          foraging, sharing resources via gifts and a common pool, and making
          use of government structures. The game is broken down into turns and
          seasons.{' '}
        </p>
        <p>
          A turn consists of different organisations running inter-island tasks.
          First the inter-island governmental organisation (IIGO) handles the
          common pool of resources and the rules of play. Then, the islands
          forage by going fishing or deer hunting. Next, the inter-island
          forecast organisation (IIFO) allows islands to share foraging and
          disaster information. This is followed by the inter-island trade
          organisation (IITO), which runs gift and information exchanges between
          islands. The runnings of the IIGO, IIFO and IITO are led by the
          islands in the roles of president, judge and speaker. Islands are
          elected into these positions of power. Finally, taxes and the cost of
          living are deducted from islands and the turn ends.{' '}
        </p>
        <p>
          A season is made up of one or more turns, and ends when a disaster
          happens. The disaster depletes resources from the common pool and the
          islands proportional to its magnitude and epicentre location.{' '}
        </p>
        <p>
          The aim of the game is for as many islands as possibe to survive for
          as long as possible.
        </p>
        <p>
          Through this game, we explore different aspects of SOMAS, such as long
          and short-term collective risk dilemmas, as well as social dilemmas.
          Each team created their own agent island to try to best survive the
          game.
        </p>
        <div style={{ paddingTop: 5 }}>
          <h2>The Agents</h2>
          <AgentTable />
        </div>
        <div style={{ paddingTop: 20 }}>
          <h2>How to play</h2>
          <p>
            Click New Run along the top navigation bar. Here, you can choose to
            run the game with the default flags, or you can customize the game
            by changing the flags. Click to run the game. From here, you can
            either download the outputs or click Visualise and explore different
            diagrams showing the progression of the game ran. The CI Output
            (continuous integration output) tab gives the automated outputs of
            the run generated from our code. The code for this can be seen{' '}
            <a href="https://github.com/SOMAS2020/SOMAS2020">here.</a>
          </p>
          <p>
            The Visualisation tab will display the results of your own runs and
            it will persist your visualisations until you visualise a new game.
          </p>
        </div>
      </div>
    </div>
  )
}

export default About
