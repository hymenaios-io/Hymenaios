import { FC, JSX, memo } from 'react';
import { OverlayTrigger, Tooltip } from 'react-bootstrap';

import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faQuestionCircle } from '@fortawesome/free-solid-svg-icons';

interface Props {
	text: string | JSX.Element;
	placement?: 'top' | 'right' | 'bottom' | 'left';
}

/**
 * A tooltip inside a question mark icon.
 *
 * @param text - The text to display in the tooltip.
 * @param placement - The placement of the tooltip.
 * @returns A hoverable tooltip inside a question mark icon.
 */
const HelpTooltip: FC<Props> = ({ text, placement = 'top' }) => (
	<OverlayTrigger
		placement={placement}
		delay={{ show: 500, hide: 500 }}
		overlay={<Tooltip id="help-tooltip">{text}</Tooltip>}
	>
		<FontAwesomeIcon
			icon={faQuestionCircle}
			style={{
				paddingLeft: '0.25em',
				height: '0.75em',
			}}
		/>
	</OverlayTrigger>
);

export default memo(HelpTooltip);
