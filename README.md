# durago
Implementation of ISO-8601 duration parsing

Inspired by [sosodev's implementation](https://github.com/sosodev/duration)

# Purposes
I needed this parser for work purposes.  
sosodev's solution was not suitable for me, because it did not handle the parsing errors I needed and allows you to specify much more than I needed in the values.

# Accuracy and values
The package provides stronger validation than sosodev's solution, emphasizing restricting values ​​to integers except seconds.  
Duration values ​​with years or months will be converted with a slight inaccuracy, as the values ​​vary.
Similar to sosodev's solution, `2.628e+15` nanoseconds for a month and `3.154e+16` nanoseconds for a year are used.

# Future improvements
- [ ] Add conversion from `time.Duration` to `durago.Duration`
- [ ] Add formatting to `ISO-8601 duration`
- [ ] Add marshaller and unmarshaller interface

# Usage
// TODO
