import {Footer} from '../components/Footer';
import {connect} from 'react-redux';

const mapStateToProps = (state) =>
    state.apiinfo.version ? state.apiinfo : {version:""};

export const FooterContainer = connect(mapStateToProps)(Footer);

