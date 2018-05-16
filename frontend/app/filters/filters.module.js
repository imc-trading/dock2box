angular
  .module("app.filters", [])
  .filter("booleanSearch", booleanSearch)
  .filter("range", range)
  .filter("joinBy", joinBy)
  .filter("replace", replace)
  .filter("fmtDate", fmtDate)
  .filter("fmtDuration", fmtDuration)
  .filter("fmtField", fmtField);

function booleanSearch() {
  return function(rows, search, fields) {
    if (typeof search === "undefined") {
      return rows;
    }

    var exprs = search.split(/\s+/);

    var matchesLength = 0;
    var includesLength = 0;
    var excludesLength = 0;
    for (var i = 0; i < exprs.length; i++) {
      var expr = exprs[i].replace("+", "").replace("-", "");
      var prefix = exprs[i].charAt(0);

      if (expr == "") {
        continue;
      }

      switch (prefix) {
        case "+":
          includesLength++;
          break;
        case "-":
          excludesLength++;
          break;
        default:
          matchesLength++;
      }
    }

    var results = [];
    for (var i = 0; i < rows.length; i++) {
      var row = rows[i];
      var found = false;

      rows[i].matches = 0;
      rows[i].includes = 0;
      rows[i].excludes = 0;

      for (var j = 0; j < exprs.length; j++) {
        var expr = exprs[j].replace("+", "").replace("-", "");
        var prefix = exprs[j].charAt(0);
        var match = false;
        var include = false;
        var exclude = false;

        if (expr == "") {
          continue;
        }

        for (var k in row) {
          if (typeof fields !== "undefined") {
            if (fields.indexOf(k) < 0) {
              continue;
            }
          }
          var val = row[k].toString();

          if (val.toLowerCase().indexOf(expr.toLowerCase()) !== -1) {
            switch (prefix) {
              case "+":
                include = true;
                break;
              case "-":
                exclude = true;
                break;
              default:
                match = true;
            }
          }
        }

        if (match) {
          rows[i].matches++;
        }

        if (include) {
          rows[i].includes++;
        }

        if (exclude) {
          rows[i].excludes++;
        }
      }

      if (
        (rows[i].matches > 0 || matchesLength == 0) &&
        rows[i].includes == includesLength &&
        rows[i].excludes == 0
      ) {
        results.push(row);
      }
    }

    return results;
  };
}

function range() {
  return function(input, start, total) {
    for (var i = start; i <= total; i++) {
      input.push(i);
    }
    return input;
  };
}

function joinBy() {
  return function(input, delimiter) {
    return (input || []).join(delimiter || ",");
  };
}

function replace() {
  return function(input, restr, newstr) {
    var re = new RegExp(restr, "g");
    return input.replace(re, newstr);
  };
}

function fmtDate() {
  return function(input) {
    var d = new Date(input);
    if (isNaN(d.getMonth())) {
      return "";
    }

    return (
      d.getFullYear() +
      "-" +
      ("0" + (d.getMonth() + 1)).slice(-2) +
      "-" +
      ("0" + d.getDate()).slice(-2) +
      " " +
      ("0" + d.getHours()).slice(-2) +
      ":" +
      ("0" + d.getMinutes()).slice(-2)
    );
  };
}

function fmtDuration() {
  return function(input) {
    if (isNaN(input)) {
      return;
    }

    dur = new Date(Math.abs(input) / 1000000);
    str = "";

    if (dur.getUTCHours() != 0) {
      str += dur.getUTCHours() + " h ";
    }

    if (dur.getUTCMinutes() != 0) {
      str += dur.getUTCMinutes() + " m ";
    }

    str += dur.getUTCSeconds() + " s";
    return str;
  };
}

function fmtField() {
  return function(val, kind) {
    switch (kind) {
      case "datetime":
        return fmtDate()(val);
        break;
      case "duration":
        return fmtDuration()(val);
      case "number":
        if (typeof val === "undefined") {
          return val;
        }
        return val.toFixed(2);
      default:
        return val;
    }
  };
}
