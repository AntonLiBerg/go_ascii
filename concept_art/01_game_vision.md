# Game Vision

## Working Concept

A turn-based, top-down command game in which the player captains a ship from its bridge. The captain physically moves between instruments to gather information, steer the ship, manage its systems, and issue orders to the crew. Every movement and interaction advances time, so each situation becomes a puzzle about what to inspect, what to act on, and in which order.

## Player Fantasy

The player should feel like a captain making difficult decisions with incomplete information. They do not perform every job personally. They decide what deserves attention, place crew where they are needed, and keep the whole ship working toward its objective.

## Design Pillars

1. **The bridge is the interface.** Information and actions are accessed through instruments in the physical game world, not through an all-purpose menu.
2. **Information costs time.** Inspecting a room or scanning a hazard uses a turn that could have been spent acting.
3. **Orders are delayed.** Crew members need time to reach a room and complete a task.
4. **Problems overlap.** Navigation, external hazards, ship damage, and crew needs compete for the captain's attention.
5. **Several plans can work.** Encounters should reward judgment rather than require one hidden sequence of actions.

## Core Loop

1. Notice alarms, reports, and visible changes on the bridge.
2. Move to an instrument.
3. Inspect information or perform an action.
4. Advance the ship, crew, and active problems by one turn.
5. Adapt the plan until the immediate encounter is resolved.
6. Continue along the larger voyage with the consequences of those decisions.

## Turn Model

On a turn, the captain does one of the following:

- Move one tile orthogonally.
- Interact with an adjacent instrument.
- Wait.

After the captain acts, the world advances once:

- The ship continues on its current course and speed.
- Crew members move or work on their assigned orders.
- Hazards approach or change.
- Damage and other unresolved problems progress.

Opening an instrument and choosing its action are one interaction. Reading ordinary interface text does not consume additional turns.

## Feedback Rules

- The captain always sees their immediate surroundings and obvious bridge alarms.
- Detailed information is only current when obtained from the relevant instrument.
- Previously gathered information remains visible but is marked with the turn when it was observed.
- Orders report when they are acknowledged, started, completed, or interrupted.
- Consequences should be predictable from information the player could have obtained.

## Initial Scope

The first playable version should contain:

- One bridge map.
- Five or six instruments.
- A small crew represented by names, location, and current order.
- One encounter combining an external hazard with an internal problem.
- Clear success, failure, and partial-success outcomes.

The first version does not need combat, inventory, character dialogue, procedural generation, or a full voyage campaign.

## Questions To Resolve Through Play

- Does walking between instruments create decisions or only delay?
- How much information can remain visible before instruments lose their purpose?
- How many simultaneous problems can the player understand comfortably?
- Should changing speed alter how often the outside world advances?
- What consequences should persist between encounters?
