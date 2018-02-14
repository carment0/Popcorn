/**
 * @copyright Popcorn, 2018
 * @author Calvin Feng
 */

// Thirdparty imports
import React from 'react';
import PropTypes from 'prop-types';
import Slider from 'rc-slider';
import { Button, OverlayTrigger, Popover } from 'react-bootstrap';

// Style imports
import './popularity_selector.scss';


class PopularitySelector extends React.Component {
  state = {
    percentile: this.props.moviePopularityPercentile
  };

  static propTypes = {
    moviePopularityPercentile: PropTypes.number.isRequired,
    disabled: PropTypes.bool.isRequired,
    dispatchSetMoviePopularityPercentile: PropTypes.func.isRequired
  };

  handleChange = (value) => {
    this.setState({ percentile: value });

    if (this.dispatchTimer) {
      clearTimeout(this.dispatchTimer);
    }

    this.dispatchTimer = setTimeout(() => {
      this.props.dispatchSetMoviePopularityPercentile(this.state.percentile);
    }, 1000);
  };

  render() {
    let title;

    switch (this.state.percentile) {
      case 0:
        title = 'Any movie';
        break;
      case 20:
        title = '> 20th percentile';
        break;
      case 50:
        title = '> 50th percentile';
        break;
      case 80:
        title = '> 80th percentile';
        break;
      case 100:
        title = 'Most highly reviewed movies';
        break;
    }

    const marks = {
      0: 'Any',
      20: '> 20th',
      50: '> 50th',
      80: '> 80th',
      100: 'Top'
    };

    const popover = (
      <Popover id="popover-positioned-bottom" title={title} className="popularity-selector">
        <div className="instruction">
          <p>Slide to filter by reviews</p>
        </div>
        <div className="slider">
          <Slider
            disabled={this.props.disabled}
            min={0}
            marks={marks}
            step={null}
            onChange={this.handleChange}
            defaultValue={this.state.percentile} />
        </div>
      </Popover>
    );

    return (
      <OverlayTrigger trigger="click" placement="bottom" overlay={popover}>
        <Button className="popularity-selector-trigger">Filter by Reviews</Button>
      </OverlayTrigger>
    );
  }
}

export default PopularitySelector;
