# IIGO Voteable Rules

This is a list of the rules available to agents, when it is worth checking if they are in play, and where to find the values you should check. Hopefully this also helps you consider which rules your agent will want to see in play. You can find the definitions in English of all the rules in the design spec and the definitions in GoLang in *globalrulesregistrations.go*

Rule names that include |rolename| indicate that there are 3 rules, one for each role



## Budget & Salary

**Name**:

- president_over_budget
- speaker_over_budget
- judge_over_budget

Logic: |rolename|LeftoverBudget >= 0

When to check: Every function where the agent decides to perform an IIGO action with a cost. You can find the costs of IIGO budget in the exposed config. You can find the remaining budget in the client gamestate. The functions affected by costs are:

- GetRuleForSpeaker
- BroadcastTaxation
- ReplyAllocationRequests
- RequestAllocationRequest
- RequestRuleProposal
- AppointNextSpeaker
- InspectHistory
- HistoricalRetribution
- InspectBallot
- InspectAllocation
- AppointNextPresident
- SetVotingResult
- SetRuleToVote
- AnnounceVotingResult
- UpdateRulesAction
- AppointNextJudge

Note: In situations of two rules clashing (you are obliged to do an action & you are not permitted to do that action) this rule (the no permission rule) takes priority (the other will not be checked). 

------

**Name**: increment_budget_|rolename|

Variables: |rolename|BudgetIncrement

Logic: |rolename|BudgetIncrement - const == 0

When to check: the rule is never checked. 

This rule is **mutable**. It is used to set how much the budget for a role is incremented each turn. The amount roles have in budget can be accessed in gamestate. 

------

**Name**: salary_cycle_|rolename|

Variables:|rolename|_Payment

Logic: |rolename|Payment - const == 0

When to check: In |rolename|.Pay|rolename|()

This rule is **mutable**. By changing the constant islands can adjust the rule how much, if payment has happened, should be paid in salary. 

------

**Name**: salary_paid_|rolename|

Variables: |rolename|Paid

Logic: |rolename|Paid - 1 == 0

When to check: In |rolename|.Pay|rolename|()

return ActionTaken bool must be true

------





## President

**Name**: allocations_made_rule

Variables: AllocationMade

Logic: AllocationMade - 1 ==0

When to check: *president.EvaluateAllocationRequests()* 

Rule states you are obliged to perform the action (PresidentReturnContent.ActionTaken  = true).

------

**Name**: rule_chosen_from_proposal_list

Variables: RuleChosenFromProposalList

Logic: RuleChosenFromProposalList- 1 == 0

When to check: President.PickRuleToVote()

------

**Name**: obl_to_propose_rule_if_some_are_given

Variables: IslandsProposedRules, RuleSelected

Logic: IslandsProposedRules - RuleSelected == 0

When to check: President.PickRuleToVote()

Must return ActionTaken and none-empty rule for the rule to be passed.

------





## Speaker

**Name**: vote_called_rule

Variables: RuleSelected, VoteCalled

Logic: RuleSelected - VoteCalled == 0

When to check: Speaker.DecideAgenda() & Speaker.DecideVote()

If the ruleMatrix passed to Speaker.DecideAgenda() is empty, the President did not select a rule. The rule states the Speaker is not allowed to perform these actions in this case. *VoteCalled*  depends on the *ActionTaken* boolean returned in Speaker.DecideVote(). 

------

**Name**: vote_result_rule

Variables: VoteResultAnnounced, VoteCalled

Logic: VoteResultAnnounced - VoteCalled == 0

When to check: Speaker.DecideVote() & Speaker.DecideAnnouncement()

If a vote is called you are obliged to announce it. If a vote is not called you are not permitted to announce it.

------

**Name**: islands_allowed_to_vote_rule

Variables: AllIslandsAllowedToVote

Logic: AllIslandsAllowedToVote - 1 == 0

When to check: Speaker.DecideVote()

The function is passed IDs of all alive islands. The vote must be called for all alive islands participating (aka change nothing). 

------

**Name**: rule_to_vote_on_rule

Variables: SpeakerProposedPresidentRule

Logic: SpeakerProposedPresidentRule - 1 == 0

When to check: Speaker.DecideAgenda() & Speaker.DecideVote()

In Speaker.DecideAgenda() the island is passed the rule the president set out for voting. Speaker.DecideVote() is passed the decided agenda rule. The Speaker.DecideVote() returned rule must equal the president decided one first given to the island in Speaker.DecideAgenda() for this rule to pass.

------

**Name**: announcement_matches_vote

Variables: AnnouncementRuleMatchesVote, AnnouncementResultMatchesVote

Logic: Both need to be true

When to check: Speaker.DecideAnnouncement()

The announced Rule and corresponding result must match what was voted on.





## Judge

**Name**: judge_inspection_rule

Variables: JudgeInspectionPerformed

Logic: JudgeInspectionPerformed- 1 ==0

When to check: *Judge.InspectHistory*

Rule states you are obliged to perform the action (return boolean  = true). Note: Inspect history is left completely up to islands, one of the few instances where you can cheat as a role and nobody can find out (such is the nature of time constrained implementation). 

------

**Name**: judge_historical_retribution_permission

Variables: JudgeHistoricalRetributionPerformed

Logic: JudgeHistoricalRetributionPerformed == 0

When to check: Judge.HistoricalRetributionEnabled()

If rule is enabled the Judge is not permitted to do historical retribution (check island actions further back than only last turn). Must return false

------





## Island

**Name**: island_must_report_private_resource

Variables: HasIslandReportPrivateResources

Logic: HasIslandReportPrivateResources - 1 == 0

When to check: ResourceReport()

Rule states you must report something (return bool = true)

------

**Name**: island_must_report_actual_private_resource

Variables: IslandActualPrivateResources, IslandReportedPrivateResources,

Logic: IslandActualPrivateResources - IslandReportedPrivateResources == 0

When to check: ResourceReport()

Rule states if you report, you must report your true amount

------





## Elections & Monitoring

**Name**: iigo_monitor_rule_permission_1

Variables: MonitorRoleDecideToMonitor, MonitorRoleAnnounce

Logic: MonitorRoleDecideToMonitor - MonitorRoleAnnounce == 0

When to check: |rolename|.MonitorIIGORole(),  |rolename|.DecideIIGOMonitoringAnnouncement()

If monitoring has occurred, the rule states you should return _, true for DecideIIGOMonitoringAnnouncement()

------

**Name**: iigo_monitor_rule_permission_2

Variables: MonitorRoleEvalResult, MonitorRoleEvalResultDecide

Logic: MonitorRoleEvalResult - MonitorRoleEvalResultDecide == 0

When to check: |rolename|.DecideIIGOMonitoringAnnouncement()

Rules states you must not change result.

------

**Name**: roles_must_hold_election

Variables: TermEnded, ElectionHeld

Logic: TermEnded - ElectionHeld == 0

When to check: |rolename1|.Call|rolename2|Election

To check if a roles term has ended check clientgamestate.IIGOTurnsInPower[|rolename2|] > gameConf.IIGOTermLengths[|rolename2|]

------

**Name**: must_appoint_elected_island

Variables: AppointmentMatchesVote

Logic: AppointmentMatchesVote - 1 == 0

When to check: |rolename1|.DecideNext|rolename2|()

Unlike voting for rules, the announcement is always made, but the islands still have the power of deciding what to announce. Hence this rule states that the announcement must result the result of the election (aka return the clientID given)