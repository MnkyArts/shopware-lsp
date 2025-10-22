<?php declare(strict_types=1);

namespace Swag\CustomProduct\Core\Content\CustomProduct;

use Shopware\Core\Framework\DataAbstractionLayer\EntityDefinition;
use Shopware\Core\Framework\DataAbstractionLayer\FieldCollection;

class CustomProductDefinition extends EntityDefinition
{
    public const ENTITY_NAME = 'custom_product';

    public function getEntityName(): string
    {
        return self::ENTITY_NAME;
    }

    protected function defineFields(): FieldCollection
    {
        return new FieldCollection([
            // fields here
        ]);
    }
}
