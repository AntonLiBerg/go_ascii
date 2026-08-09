# Scenarios

Each scenario is a directory containing `map.txt`, `content.txt`, and `ui.txt`.

## Map

Declare one or more named rooms. A room may list multiple entity groups.

```text
===bridge: base,terminals_bridge
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
- `group1,group2` after a room name combines those entity groups.
- `portal:<marker>,<room>,<marker>` connects two markers for movement in both directions.
- `terminal:<entity marker>,<room>` opens another room when that neighboring entity has `interactable:terminal`.
- `selectableorder:<marker>,...` defines the selection order for controls in a room.

A terminal target room is rendered as literal ASCII, so text and symbols do not need entity definitions. Its `exit` input closes the terminal without moving the player.

## UI

`ui.txt` defines the vertical draw order and static UI sections:

```text
===layout
room
infobox
===infobox
+----------------+
|                |
+----------------+
```

`room` draws the active room. Other names draw their matching section below or above it in layout order. The room is centered within the widest UI section when it fits.

## Content

Entity groups map characters to entities. A room uses the union of its groups:

```text
===base
.:floor
- pos
- ascii:.

@:player
- pos
- ascii:o
- player
```

Unused entity definitions and groups are allowed.

Terminal rooms may use the built-in `terminal` group. Their map text is rendered literally.

Use `ascii:SPACE` when an entity should render as a blank space.

An interactive entity declares one interaction type with the `interactable` component:

```text
D:door
- pos
- ascii:+
- impassable
- interactable:door
```

The supported interaction types are `door` and `terminal`.

A selectable entity declares its unfocused, focused, and selected ASCII followed by its target entity:

```text
1:selectposition
- pos
- ascii:o
- selectable:o,ö,^,facing
```

The first entry in a room's `selectableorder` is focused whenever that room is entered. The selected ASCII is shown while editing that focused control.

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

The active room supplies the movement and interaction profile. Terminal rooms normally define only `quitgame` and `exit`.
