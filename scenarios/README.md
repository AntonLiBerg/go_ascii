# Scenarios

Each scenario is a directory containing `map.txt` and `content.txt`.

## Map

Declare one or more named rooms. Every room needs a ground entity and an input profile.

```text
===bridge
-----0-----
|    @    |
-----------
features
- ground:floor
- inputprofile:topdown
- portal:0,comms,0
```

Available room features:

- `ground:<entity>` places that entity below every map cell.
- `inputprofile:<name>` selects a profile from `content.txt` while viewing the room.
- `portal:<marker>,<room>,<marker>` connects two markers for movement in both directions.
- `terminal:<entity marker>,<room>` opens another room when that neighboring entity has `interactable:commandtable`.

A terminal target room is rendered as literal ASCII, so text and symbols do not need entity definitions. Its `exit` input closes the terminal without moving the player.

## Content

Entity headers map one character to an entity name, followed by components:

```text
===ENTITY
.:floor
- pos
- ascii:.

@:player
- pos
- ascii:o
- player
```

Use `ascii:SPACE` when an entity should render as a blank space.

An interactive entity declares one interaction type with the `interactable` component:

```text
D:door
- pos
- ascii:+
- impassable
- interactable:door
```

The supported interaction types are `door`, `helm`, and `commandtable`.

Define named input profiles after the entities:

```text
===INPUTPROFILE
topdown
- quitgame=q
- interact=e
- moveup=w
- moveleft=a
- movedown=s
- moveright=d

terminal_scan
- quitgame=q
- exit=e
```
