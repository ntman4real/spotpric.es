var uniqueZones = function(prices) {
  var zones = [];
  prices.forEach(function(p) {
    if (zones.indexOf(p.Zone) == -1) zones.push(p.Zone);
  });
  return zones.sort();
};

var SpotHeader = React.createClass({
  render: function() {
    return (
      <thead>
        <tr>
          <th>Instance Type</th>
          {this.props.zones.map(function(z) { return <th>{z}</th>; })}
        </tr>
      </thead>
    );
  }
});

var SpotPrice = React.createClass({
  getInitialState: function() {
    return { changed: false, changeDown: false };
  },

  componentWillReceiveProps: function(nextProps) {
    if (nextProps.price == null) {
      this.setState({ changed: false, changeDown: false });
      return;
    }

    if (this.props.price == null || (nextProps.price.Price != this.props.price.Price)) {
      var newState = { changed: true, changeDown: false };

      if (this.props.price != null && this.props.price.Price > nextProps.price.Price) {
        newState.changeDown = true;
      }

      this.setState(newState);

      setTimeout(function() {
        // just reset everything
        this.setState({ changed: false, changeDown: false });
      }.bind(this), 15000);
    } else {
      this.setState({ changed: false, changeDown: false });
    }
  },

  render: function() {
    var shouldHighlight = this.state.changed && this.state.changeDown;

    return (
      <td className={this.state.changed ? (this.state.changeDown ? 'highlight' : 'highlight-red') : ''}>
        <span className={this.state.changed ? 'animated bounceIn' : ''}>{this.props.price == null ? ' ' : this.props.price.Price.toFixed(3)}</span>
      </td>
    );
  }
});

var SpotRow = React.createClass({
  render: function() {
    return (
      <tr>
        <td>{this.props.group[0].InstanceType}</td>
        {this.props.zones.map(function(zone) { 
          var fz = _.find(this.props.group, function(g) { return g.Zone == zone; })
          return <SpotPrice price={fz} />;
        }.bind(this))}
      </tr>
    );
  }
});

var SpotPrices = React.createClass({
  getInitialState: function() {
    return { 
      prices: window.__initialPrices 
    };
  },

  componentDidMount: function() {
    var conn = new ReconnectingWebSocket("ws://"+window.location.host+"/ws");
    conn.onmessage = function(evt) {
      var prices = JSON.parse(evt.data);
      this.setState({ prices: prices });
    }.bind(this);
  },

  render: function() {
    var doc;

    if (this.state.prices == null) {
      doc = <h2>Hold Your Horses!!!</h2>;
    } else {
      var zones = uniqueZones(this.state.prices);
      var groupedByType = _.groupBy(this.state.prices, function(p) { return p.InstanceType });

      doc = (
        <table className="fix">
          <SpotHeader zones={zones} />
          {_.keys(groupedByType).sort().map(function(key) { return <SpotRow zones={zones} group={groupedByType[key]} />; })}
        </table>
      );
    }

    return (
      <div>
        {doc}
      </div>
    );
  }
});

React.render(
  <SpotPrices />,
  document.getElementById('table')
);

