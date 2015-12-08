# picsv

Convert pictures to svg animation written in golang

## Usage

```
picsv -o output.svg \
-pics_1 background.gif \
-pics_2 front_1.gif,front_2.gif,front_3.gif,front_4.gif \
-pics_3 bird_1.gif,bird_2.gif,bird_3.gif,bird_4.gif \
-begin_3 3.5 \
-repeat_3 1 \
-pics_4 snow_1.gif,snow_2.gif,snow_3.gif,snow_4.gif \
-begin_4 6.5 \
-repeat_4 3 \
```

## Installation
```
go get github.com/high5/picsv
```

## Options

Numeric part of the option is meant a layer.

pics_*: Image specified in the comma-separated
```
picsv -o output.svg -pics_1 front_1.gif,front_2.gif,front_3.gif
```

o: output
```
picsv -o output.svg -pics_1 front_1.gif,front_2.gif,front_3.gif
```


begin_*: This option defines when an animation should begin (for example, write '4'  begin after 4.5 seconds).
```
picsv -o output.svg -pics_1 front_1.gif -begin_1 4
```

repeat_*: Repeat count.
```
picsv -o output.svg -pics_1 front_1.gif -repeat_1 3
```

dur_*: The duration of the animation (for example, write '5' for 5 seconds).
```
picsv -o output.svg -pics_1 front_1.gif -begin_1 4
```

## License

MIT
