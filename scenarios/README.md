# Scenarios

Each scenario is a directory containing:

- `map.txt`: rooms, map cells, and room features.
- `content.txt`: entity groups, components, and input profiles.
- `ui.txt`: the layout and static UI sections.

## Map

Declare one or more named rooms with an `===<room>` header. Add a
comma-separated list after `:` to restrict the room to the union of those
entity groups.

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

- `ground:<entity>` places that entity below every map cell. Required.
- `inputprofile:<name>` selects a profile from `content.txt`. Required.
- `portal:<marker>,<room>,<marker>` connects two markers for movement in both directions.
- `terminal:<entity marker>,<room>` opens another room when that neighboring entity has `interactable:terminal`.
- `selectableorder:<marker>,...` defines the selection order for controls in a room.

Feature entries follow a `features` line. A `//` comment may follow a feature.
Portal and terminal markers must occur exactly once in their source room. A
portal is bidirectional and should be declared only once.

Spaces contain only the room's ground entity. In normal rooms, every other map
character must belong to one of the room's groups and have an entity
definition. Rooms without a group list may use any defined entity.

A terminal target room renders unknown characters as literal ASCII, so its
text and symbols do not need entity definitions. Known characters in its
declared groups still create normal entities. The built-in `terminal` group
can be used for a fully literal terminal room without defining that group in
`content.txt`.

## UI

The `layout` section defines the vertical draw order. Every other section
contains static UI lines:

```text
===layout
room
infobox
===infobox
+----------------+
|                |
+----------------+
```

`room` draws the active room. Other names draw their matching section in layout
order. Layout names must be unique, UI section names must be unique, and the
`layout` section is required. The room is centered within the widest UI
section when it fits.

## Content

Entity groups start with an `===<group>` header. Entity headers use
`<single character>:<unique name>`, followed by component lines:

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

Entity keys and names must be unique across all groups. Components for an
entity must also be unique. Unused entity definitions and groups are allowed.

Use `ascii:SPACE` when an entity should render as a blank space.

Supported components:

- `pos` gives the entity a position.
- `ascii:<character>` defines its rendered character.
- `impassable` blocks movement.
- `player` marks the player entity.
- `interactable:door` makes an impassable neighboring entity openable.
- `interactable:terminal` opens the room configured by a `terminal` map feature.
- `controlnumber:<minimum>,<maximum>` defines a numeric control value.
- `selectable:<unfocused>,<focused>,<selected>,<target entity>` controls another entity.

An interactive entity declares one interaction type with the `interactable` component:

```text
D:door
- pos
- ascii:+
- impassable
- interactable:door
```

A selectable entity declares its unfocused, focused, and selected ASCII followed by its target entity:

```text
1:selectposition
- pos
- ascii:o
- selectable:o,ö,^,facing
```

The target entity must exist and contain a compatible control component. Every
marker in `selectableorder` must resolve to a selectable entity in that room.
The first entry is focused when the room is entered, and the selected ASCII is
shown while editing the focused control.

Define named input profiles after the entities:

```text
===INPUTPROFILE
topdown
- profiletype=none
- quitgame=q
- interact=e
- moveup=w
- moveleft=a
- movedown=s
- moveright=d

terminal_scan
- profiletype=terminal
- quitgame=q
- exit=e
```

Each profile name must be unique, as must each action within a profile.
`profiletype` accepts `none`, `terminal`, or `control`. These short names are
normalized internally. Supported actions are:

- `quitgame`
- `exit`
- `moveup`, `moveleft`, `movedown`, and `moveright`
- `interact`
- `moveselectnext` and `moveselectprev`
- `select`

The active room selects its input profile. Movement rooms generally use
`profiletype=none`. Literal terminal rooms can use `profiletype=terminal`, and
rooms with selectable controls use `profiletype=control`.
