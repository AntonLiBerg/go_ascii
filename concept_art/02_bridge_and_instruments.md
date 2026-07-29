# Bridge And Instruments

## Purpose

The bridge is both the play space and the game's interface. Its layout should force useful tradeoffs without making routine actions tedious. Instruments used together often should be near each other; instruments representing competing responsibilities can be farther apart.

## Concept Layout

```text
#########################
# N                   T #
#                       #
#     C           I     #
#                       #
# E         @         H #
#                       #
########### D ###########
```

Legend:

- `@` captain's starting position
- `H` helm
- `T` tactical display
- `N` navigation chart
- `I` internal monitor
- `C` communications station
- `E` engineering console
- `D` bridge door

This is a relationship diagram, not a final map. Walking distances must be tested in play.

## Instruments

| Instrument | Information | Actions | Intended decision |
| --- | --- | --- | --- |
| Helm | Current heading and speed | Change heading or speed | How should the ship move right now? |
| Tactical display | Nearby objects, ships, and hazards | Focus a scan on one contact | What presents an immediate external danger? |
| Navigation chart | Larger route, destination, and known regions | Set a waypoint | What course supports the voyage rather than only the current crisis? |
| Internal monitor | Last known room, system, and crew status | Inspect one room in detail | Where is something going wrong inside the ship? |
| Communications station | Crew locations, orders, and availability | Assign, cancel, or reprioritize one order | Who should handle a problem, and when? |
| Engineering console | Power allocation and major system condition | Shift limited power between systems | Which capability matters most at this moment? |

## Instrument Rules

- The captain must stand next to an instrument to use it.
- One interaction provides one focused piece of information or performs one action.
- An instrument only shows information appropriate to its role.
- Information records the turn on which it was observed and may become stale.
- Actions take effect before the world advances unless the action states that it is delayed.
- An instrument should not duplicate another instrument's main purpose.

## Crew Orders

Crew members have only three required properties for the first prototype:

- Current room.
- Current order.
- Number of turns until the next order step is complete.

An order follows this sequence:

1. The captain assigns a crew member at the communications station.
2. The crew member travels toward the target room each world turn.
3. The crew member works for the stated number of turns.
4. The result is reported to the bridge.

Crew expertise, health, relationships, and autonomous priorities should be added only if the basic ordering puzzle proves interesting.

## Interaction Presentation

Instrument screens should be short and action-oriented. A screen should show:

- The instrument name.
- The current turn.
- Current or last observed information.
- Available actions and their consequences when known.
- A clear way to close the instrument without spending a turn.

Moving the selection inside an instrument does not spend turns. Confirming an inspection or action spends one turn.
