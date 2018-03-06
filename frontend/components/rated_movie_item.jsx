/**
 * @copyright Popcorn, 2018
 * @author Calvin Feng
 */

// Thirdparty imports
import React from 'react';
import PropTypes from 'prop-types';
import Spinner from 'spin';
import { OverlayTrigger, Tooltip } from 'react-bootstrap';

// Style
import './rated_movie_item.scss';


class RatedMovieItem extends React.Component {
  static propTypes = {
    imdbId: PropTypes.string.isRequired,
    movieId: PropTypes.number.isRequired,
    movieRating: PropTypes.number.isRequired,
    dispatchMovieDetailFetch: PropTypes.func.isRequired,
    movieDetail: PropTypes.object,
  }

  static defaultProps = {
    movieDetail: {
      isDefaultProp: true,
      title: 'Loading...',
      year: '....',
      plot: 'Loading...'
    }
  }

  /**
   * Initiate fetching metadata of the current movie and spinning the loading bar so user knows that it is fetching
   * something.
   */
  fetchMovieDetail = () => {
    if (this.props.movieDetail.isDefaultProp) {
      this.spinner = new Spinner().spin();
      document.getElementById(this.props.movieId).appendChild(this.spinner.el);
      this.props.dispatchMovieDetailFetch(this.props.imdbId);
    }
  };

  /**
   * @return {React.Element}
   */
  get poster() {
    if (this.props.movieDetail.poster) {
      return (
        <section className="poster">
          <img id={this.props.movieId} src={this.props.movieDetail.poster} alt="poster" />
        </section>
      );
    }

    return <section className="poster" />;
  }

  componentWillReceiveProps(nextProps) {
    if (!nextProps.movieDetail.isDefaultProp && this.spinner) {
      this.spinner.stop();
    }
  }

  componentDidMount() {
    if (this.props.movieDetail.isDefaultProp) {
      this.fetchMovieDetail();
    }
  }

  render() {
    const tooltip = (
      <Tooltip id="tooltip">
        Your rating: <strong>{this.props.movieRating}</strong>
      </Tooltip>
    );

    return (
      <section className="rated movie-item" id={this.props.movieId}>
        <OverlayTrigger placement="bottom" overlay={tooltip}>
          {this.poster}
        </OverlayTrigger>
      </section>
    );
  }
}

export default RatedMovieItem;
