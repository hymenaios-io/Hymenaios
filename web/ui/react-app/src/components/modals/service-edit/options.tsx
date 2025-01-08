import { FC, memo, useMemo } from 'react';
import { FormCheck, FormText } from 'components/generic/form';

import { Accordion } from 'react-bootstrap';
import { BooleanWithDefault } from 'components/generic';
import { ServiceOptionsType } from 'types/config';
import { firstNonDefault } from 'utils';

interface Props {
	defaults?: ServiceOptionsType;
	hard_defaults?: ServiceOptionsType;
}

/**
 * The `service.options` form fields.
 *
 * @param defaults - The default values.
 * @param hard_defaults - The hard default values.
 * @returns The form fields for the `options`.
 */
const EditServiceOptions: FC<Props> = ({ defaults, hard_defaults }) => {
	const convertedDefaults = useMemo(
		() => ({
			interval: firstNonDefault(defaults?.interval, hard_defaults?.interval),
			semantic_versioning:
				defaults?.semantic_versioning ?? hard_defaults?.semantic_versioning,
		}),
		[defaults, hard_defaults],
	);

	return (
		<Accordion>
			<Accordion.Header>Options:</Accordion.Header>
			<Accordion.Body>
				<FormCheck
					name="options.active"
					label="Active"
					tooltip="Whether the service is active and checking for updates"
					size="sm"
				/>
				<FormText
					key="interval"
					name="options.interval"
					col_sm={12}
					label="Interval"
					tooltip="How often to check for both latest version and deployed version updates"
					defaultVal={convertedDefaults.interval}
				/>
				<BooleanWithDefault
					name="options.semantic_versioning"
					label="Semantic versioning"
					tooltip="Releases follow 'MAJOR.MINOR.PATCH' versioning"
					defaultValue={convertedDefaults.semantic_versioning}
				/>
			</Accordion.Body>
		</Accordion>
	);
};

export default memo(EditServiceOptions);
