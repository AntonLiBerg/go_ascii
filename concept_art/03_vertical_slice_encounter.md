# Vertical Slice Encounter

## The Broken Passage

The ship is crossing a narrow debris field when an alarm reports a problem below deck. The captain must avoid an external collision while preventing an engine-room fire from disabling the ship.

This encounter exists to test whether movement between instruments, incomplete information, and delayed crew orders create an understandable puzzle.

## Starting Situation

- The captain starts near the center of the bridge.
- The ship is moving forward at cruising speed.
- An unidentified contact is visible as a general tactical alarm.
- A general fire alarm identifies the lower deck but not the affected room.
- One engineer and one deckhand are available.
- The route has a port detour and a starboard shortcut.

The exact hazard timing is hidden until the player uses the relevant instrument.

## Discoverable Information

| Instrument | What the player can learn |
| --- | --- |
| Tactical display | Debris will cross the current heading soon. Slowing down creates more time. |
| Navigation chart | The port detour is safe but costs progress. The starboard shortcut stays on schedule but passes through denser debris. |
| Internal monitor | A fire is growing in the engine room. A crew member can contain it, but a sharp turn while it burns may damage the engine. |
| Communications station | The engineer repairs faster; the deckhand can contain the fire but cannot repair engine damage. |
| Engineering console | Power can be shifted to engines for maneuvering or to internal suppression to slow the fire. |
| Helm | The captain can slow down, turn port, turn starboard, or hold course. |

## Progression Rules

- The debris gets closer after every world turn while the ship is moving.
- Slowing the ship delays the collision but costs voyage progress.
- The fire grows after every world turn until contained by crew or suppression power.
- A crew member needs time to reach the engine room before work begins.
- A sharp maneuver while the fire is severe causes engine damage.
- The captain may act without collecting every piece of information.

Exact turn counts should be chosen during implementation and balancing. They should allow at least two successful plans and make inspecting every instrument an unsafe default.

## Outcomes

Full success:

- The ship avoids the debris.
- The fire is contained.
- The engine is not damaged.

Partial success:

- The ship survives but loses voyage progress, takes hull damage, or damages the engine.
- The consequences are explicit and would carry into a later encounter.

Failure:

- A collision destroys the ship.
- The fire disables the engine before the immediate hazard is cleared.

## Example Strategies

- Slow down first, send the engineer, then take the safe port detour.
- Shift power to suppression, use the deckhand to contain the fire, and keep the engineer available for possible damage.
- Accept minor engine damage to take the starboard shortcut and preserve voyage progress.

These are examples, not prescribed solutions.

## Playtest Questions

- Did the player understand that movement and interaction both advance time?
- Did they have enough warning to make an informed decision?
- Was stale information useful without being fully reliable?
- Did bridge distance affect the plan?
- Did crew travel time matter?
- Could the player explain why each consequence occurred?
- Were at least two approaches plausibly successful?

## Prototype Completion Criteria

The vertical slice is complete when a player can finish the encounter from beginning to end, receive one of the listed outcomes, and understand the main cause of that outcome without developer explanation.
