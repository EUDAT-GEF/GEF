import {Footer} from '../components/Footer';
import {connect} from 'react-redux';

const mapStateToProps = (state) => state.apiinfo;
export const FooterContainer = connect(mapStateToProps)(Footer);

